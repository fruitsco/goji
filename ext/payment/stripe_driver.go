package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fruitsco/goji/x/driver"
	"github.com/shopspring/decimal"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"github.com/stripe/stripe-go/v76/webhook"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StripeDriver struct {
	client *client.API
	config *StripeConfig
	log    *zap.Logger
}

var _ = Driver(&StripeDriver{})

type StripeConfig struct {
	AccessToken          string  `conf:"access_token"`
	WebhookSecretConnect *string `conf:"webhook_secret_connect"`
	WebhookSecretAccount *string `conf:"webhook_secret_account"`
	InsecureWebhooks     bool    `conf:"insecure_webhooks"`
}

type StripeDriverParams struct {
	fx.In

	Config *StripeConfig
	Log    *zap.Logger
}

func NewStripeDriverFactory(params StripeDriverParams) driver.FactoryResult[PaymentDriver, Driver] {
	return driver.NewFactory(Stripe, func() (Driver, error) {
		return NewStripeDriver(params)
	})
}

// NewStripeDriver creates a new stripe driver
func NewStripeDriver(params StripeDriverParams) (*StripeDriver, error) {
	// FIXME: here we can add more finegrained Details
	// See https://github.com/fruitsco/roma/issues/30
	// config := &stripe.BackendConfig{
	// 	LeveledLogger: &stripe.LeveledLogger{
	// 		Level: stripe.LevelInfo,
	// 	},
	// 	MaxNetworkRetries: stripe.Int64(0),
	// }
	stripe := &client.API{}
	stripe.Init(params.Config.AccessToken, nil)

	return &StripeDriver{
		client: stripe,
		config: params.Config,
		log:    params.Log.Named("payment_stripe"),
	}, nil
}

// CreateOrder creates a new order with a payment provider
func (s *StripeDriver) CreatePaymentIntent(
	ctx context.Context,
	amount int64,
	currency string,
	customerID string,
	orderID string,
) (*PaymentIntentResult, error) {
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(amount),
		Currency:      stripe.String(currency),
		TransferGroup: stripe.String(orderID),
		Customer:      stripe.String(customerID),
		// AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
		// 	Enabled: stripe.Bool(true),
		// },
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card", "paypal",
		}),
	}
	if isEUR(currency) {
		params = &stripe.PaymentIntentParams{
			Amount:        stripe.Int64(amount),
			Currency:      stripe.String(currency),
			TransferGroup: stripe.String(orderID),
			Customer:      stripe.String(customerID),
			// AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			// 	Enabled: stripe.Bool(true),
			// },
			PaymentMethodTypes: stripe.StringSlice([]string{
				"card", "paypal", "customer_balance",
			}),
		}
	}

	paymentIntent, err := s.client.PaymentIntents.New(params)

	if err != nil {
		return nil, err
	}

	return &PaymentIntentResult{
		PaymentIntentID: paymentIntent.ID,
		ClientSecret:    paymentIntent.ClientSecret,
	}, nil
}

// UpdatePaymentIntent updates an existing order with a payment provider
func (s *StripeDriver) UpdatePaymentIntent(
	ctx context.Context,
	paymentIntentID string,
	amount int64,
	currency string,
	country string,
) (*PaymentIntentResult, error) {
	params := &stripe.PaymentIntentParams{}
	setPaymentOptions(params, country, currency)

	paymentIntent, err := s.client.PaymentIntents.Update(paymentIntentID, params)

	if err != nil {
		return nil, err
	}

	return &PaymentIntentResult{
		PaymentIntentID: paymentIntent.ID,
		ClientSecret:    paymentIntent.ClientSecret,
	}, nil
}

// GetCustomerID
func (s *StripeDriver) GetCustomerForPaymentIntent(
	ctx context.Context,
	paymentIntentID string,
) (*Customer, error) {
	params := &stripe.PaymentIntentParams{}
	params.AddExpand("customer")
	paymentIntent, err := s.client.PaymentIntents.Get(paymentIntentID, params)

	if err != nil {
		return nil, err
	}

	return &Customer{
		ID:    paymentIntent.Customer.ID,
		Email: paymentIntent.Customer.Email,
	}, nil
}

// GetBillingDetailsForCharge
func (s *StripeDriver) GetBillingDetailsForCharge(
	ctx context.Context,
	chargeID string,
) (*BillingDetails, error) {
	params := &stripe.ChargeParams{}
	params.AddExpand("billing_details")
	charge, err := s.client.Charges.Get(chargeID, params)

	if err != nil {
		return nil, err
	}

	return &BillingDetails{
		ID:          charge.Customer.ID,
		Email:       charge.BillingDetails.Email,
		CountryCode: charge.BillingDetails.Address.Country,
		Name:        charge.BillingDetails.Name,
		Street:      charge.BillingDetails.Address.Line1,
		Zip:         charge.BillingDetails.Address.PostalCode,
		City:        charge.BillingDetails.Address.City,
	}, nil
}

// CreateCustomer creates a new customer with a payment provider
func (s *StripeDriver) CreateCustomer(
	ctx context.Context,
	newCustomer *Customer,
) (string, error) {
	params := &stripe.CustomerParams{}
	customer, err := s.client.Customers.New(params)
	if err != nil {
		return "", err
	}
	return customer.ID, nil
}

// GetPaymentMethode
func (s *StripeDriver) GetPaymentMethod(
	ctx context.Context,
	paymentMethodeID string,
) (string, error) {
	params := &stripe.PaymentMethodParams{}
	paymentMethode, err := s.client.PaymentMethods.Get(paymentMethodeID, params)
	if err != nil {
		return "", err
	}
	return string(paymentMethode.Type), nil
}

// UpdateCustomer
func (s *StripeDriver) UpdateCustomer(
	ctx context.Context,
	customerID string,
	customer *Customer,
) error {
	params := &stripe.CustomerParams{
		Email: stripe.String(customer.Email),
	}
	_, err := s.client.Customers.Update(customerID, params)
	return err
}

// CreateTransfer moves pending payments to their destination
func (s *StripeDriver) CreateTransfer(
	ctx context.Context,
	amount int64,
	currency string,
	recipientAccountID string,
	sourceTransactionID *string,
	transferGroupID string,
) (string, error) {
	params := &stripe.TransferParams{
		Amount:            stripe.Int64(amount),
		Currency:          stripe.String(currency),
		Destination:       stripe.String(recipientAccountID),
		SourceTransaction: sourceTransactionID,
		TransferGroup:     stripe.String(transferGroupID),
	}

	transfer, err := s.client.Transfers.New(params)

	if err != nil {
		return "", err
	}

	return transfer.ID, nil
}

func (s *StripeDriver) CreatePayout(
	ctx context.Context,
	id string,
	amount int64,
	currency string,
	recipientAccountID string,
) (*string, error) {
	params := &stripe.PayoutParams{
		Description: stripe.String(id),
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(currency),
	}
	params.SetStripeAccount(recipientAccountID)

	payout, err := s.client.Payouts.New(params)

	if err != nil {
		return nil, err
	}

	return &payout.ID, nil
}

// GetCharge gets a charge from the payment provider
func (s *StripeDriver) GetBalanceTransactionForCharge(
	ctx context.Context,
	chargeID string,
) (*BalanceTransaction, error) {
	// TODO: only get balance transaction w/o charge
	params := &stripe.ChargeParams{}
	params.AddExpand("balance_transaction")
	charge, err := s.client.Charges.Get(chargeID, params)

	if err != nil {
		return nil, err
	}

	return &BalanceTransaction{
		Amount:       decimal.NewFromInt(charge.BalanceTransaction.Amount),
		ProviderFee:  decimal.NewFromInt(charge.BalanceTransaction.Fee),
		ExchangeRate: decimal.NewFromFloat(charge.BalanceTransaction.ExchangeRate),
		Currency:     string(charge.BalanceTransaction.Currency),
	}, nil
}

// CreateTransferReversal reverses a transfer if in our application a failure occurs
func (s *StripeDriver) CreateTransferReversal(
	ctx context.Context,
	amount int64,
	transferID string,
) (*string, error) {
	params := &stripe.TransferReversalParams{
		ID:     stripe.String(transferID),
		Amount: stripe.Int64(amount),
	}

	reversal, err := s.client.TransferReversals.New(params)

	if err != nil {
		return nil, err
	}

	return &reversal.ID, nil
}

// CreateAccount creates a new account with a payment provider
func (s *StripeDriver) CreateAccount(
	ctx context.Context,
	email string,
	payoutDelay int64,
) (string, error) {
	// FIXME: Not correct if we should define default payout values here, but i guess so.
	// See https://github.com/fruitsco/roma/issues/31
	params := &stripe.AccountParams{
		Type: stripe.String("custom"),
		// BusinessType: stripe.String(string(businessType)),
		Email: stripe.String(email),
		Capabilities: &stripe.AccountCapabilitiesParams{
			CardPayments: &stripe.AccountCapabilitiesCardPaymentsParams{
				Requested: stripe.Bool(true),
			},
			Transfers: &stripe.AccountCapabilitiesTransfersParams{
				Requested: stripe.Bool(true),
			},
		},
		Settings: &stripe.AccountSettingsParams{
			Payouts: &stripe.AccountSettingsPayoutsParams{
				Schedule: &stripe.AccountSettingsPayoutsScheduleParams{

					Interval: stripe.String(string(stripe.AccountSettingsPayoutsScheduleIntervalManual)),
					// DelayDays: stripe.Int64(payoutDelay),
					// MonthlyAnchor: stripe.Int64(monthlyAnchor),
				},
			},
		},
	}

	account, err := s.client.Accounts.New(params)

	if err != nil {
		return "", err
	}

	return account.ID, nil
}

// UpdateAccount updates an account with a payment provider
func (s *StripeDriver) UpdateAccount(
	ctx context.Context,
	accountID string,
	account Account,
) error {
	params := &stripe.AccountParams{
		Email: account.Email,
		// BusinessType: stripe.String(string(account.BusinessType)),
		Country: account.CountryCode,
		Company: &stripe.AccountCompanyParams{
			Address: &stripe.AddressParams{
				City:       account.AddressCity,
				Country:    account.AddressCountryCode,
				Line1:      account.AddressStreet,
				PostalCode: account.AddressZip,
			},
			Name:  account.BusinessName,
			VATID: account.BusinessVATID,
		}}

	_, err := s.client.Accounts.Update(accountID, params)
	return err

}

// GetAccountVerificationLink gets the link to start the account verification
func (s *StripeDriver) GetAccountVerificationLink(
	ctx context.Context,
	vendorAccountID string,
	returnURL string,
	linkType string,
) (string, error) {
	params := &stripe.AccountLinkParams{
		Account: stripe.String(vendorAccountID),
		// ISSUE: We have to add the correct RefreshURL here, so that the user is redirected to the onboarding process
		RefreshURL: stripe.String(returnURL),
		ReturnURL:  stripe.String(returnURL),
		Type:       stripe.String(linkType),
		Collect:    stripe.String("eventually_due"),
	}

	accountLink, err := s.client.AccountLinks.New(params)

	if err != nil {
		return "", err
	}

	return accountLink.URL, nil
}

// GetAccountBalance gets the balance of an account
func (s *StripeDriver) GetAccountWallet(ctx context.Context, vendorAccountID string) (*Wallet, error) {
	params := &stripe.BalanceParams{}
	params.SetStripeAccount(vendorAccountID)
	balance, err := s.client.Balance.Get(params)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Available: Balance{
			amount:   balance.Available[0].Amount,
			currency: string(balance.Available[0].Currency),
		},

		Pending: Balance{
			amount:   balance.Pending[0].Amount,
			currency: string(balance.Pending[0].Currency),
		},
	}, nil
}

func (s *StripeDriver) VerifyWebhookEvent(ctx context.Context, eventData string, signature string, webhookType WebhookType) (any, error) {

	var webHookSecrect string
	if s.config.WebhookSecretConnect == nil && s.config.WebhookSecretAccount == nil && !s.config.InsecureWebhooks {
		return nil, fmt.Errorf("webhook secret not set")
	}

	if s.config.InsecureWebhooks {
		event := &stripe.Event{}
		err := json.Unmarshal([]byte(eventData), event)
		if err != nil {
			return nil, err
		}
		return event, nil
	}

	switch webhookType {
	case WebhookTypeAccount:
		webHookSecrect = *s.config.WebhookSecretAccount
	case WebhookTypeConnect:
		webHookSecrect = *s.config.WebhookSecretConnect
	}

	event, err := webhook.ConstructEvent([]byte(eventData), signature, webHookSecrect)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func isEU(country string) bool {
	return (country == string(CountryDE) || country == string(CountryFR) ||
		country == string(CountryBE) || country == string(CountryNL) ||
		country == string(CountryES) || country == string(CountryIE))
}

func isEUR(currency string) bool {
	return strings.ToLower(currency) == string(stripe.CurrencyEUR)
}
func isUSD(currency string) bool {
	return strings.ToLower(currency) == string(stripe.CurrencyUSD)
}
func isGBP(currency string) bool {
	return strings.ToLower(currency) == string(stripe.CurrencyGBP)
}

func setPaymentOptions(params *stripe.PaymentIntentParams, country string, currency string) {
	baseKey := "payment_method_options[customer_balance]"

	if isEU(country) && isEUR(currency) {
		params.AddExtra(baseKey+"[funding_type]", "bank_transfer")
		params.AddExtra(baseKey+"[bank_transfer][type]", "eu_bank_transfer")
		params.AddExtra(baseKey+"[bank_transfer][eu_bank_transfer][country]", country)
	}
	if country == string(CountryGB) && isGBP(currency) {
		params.AddExtra(baseKey+"[funding_type]", "bank_transfer")
		params.AddExtra(baseKey+"[bank_transfer][type]", "gb_bank_transfer")
	}
	if country == string(CountryUS) && isUSD(currency) {
		params.AddExtra(baseKey+"[funding_type]", "bank_transfer")
		params.AddExtra(baseKey+"[bank_transfer][type]", "us_bank_transfer")
	}

}
