package payment

import (
	"context"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	// PaymentIntent
	CreatePaymentIntent(ctx context.Context, amount int64, currency string, customerID string, orderID string) (*PaymentIntentResult, error)
	UpdatePaymentIntent(ctx context.Context, paymentIntentID string, amount int64, currency string, country string) (*PaymentIntentResult, error)
	GetPaymentMethod(ctx context.Context, paymentMethodID string) (string, error)
	// Account
	CreateAccount(ctx context.Context, email string, payoutDelay int64, account Account) (string, error)
	GetAccountVerificationLink(ctx context.Context, vendorAccountID string, returnUrl string, linkType string) (string, error)
	GetAccountWallet(ctx context.Context, vendorAccountID string) (*Wallet, error)
	UpdateAccount(ctx context.Context, accountID string, account Account) error
	// Payout
	CreatePayout(ctx context.Context, id string, amount int64, currency string, vendorAccountID string) (*string, error)
	// Transfer
	CreateTransfer(ctx context.Context, amount int64, currency string, recipientAccountID string, sourceTransactionID *string, transferGroupID string) (string, error)
	CreateTransferReversal(ctx context.Context, amount int64, transferID string) (*string, error)
	// Customer
	GetCustomerForPaymentIntent(ctx context.Context, paymentIntentID string) (*Customer, error)
	CreateCustomer(ctx context.Context, customer *Customer) (string, error)
	UpdateCustomer(ctx context.Context, customerID string, customer *Customer) error
	//Charge
	GetBillingDetailsForCharge(ctx context.Context, chargeID string) (*BillingDetails, error)
	GetBalanceTransactionForCharge(ctx context.Context, chargeID string) (*BalanceTransaction, error)
	// Verification
	VerifyWebhookEvent(ctx context.Context, eventData string, signature string, webHookType WebhookType) (any, error)
}

type PaymentParams struct {
	fx.In

	Drivers []*driver.Factory[PaymentDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Payment struct {
	drivers *driver.Pool[PaymentDriver, Driver]
	config  *Config
	log     *zap.Logger
}

func New(params PaymentParams) *Payment {
	return &Payment{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("payment"),
	}
}

func (s *Payment) resolveDriver() (Driver, error) {
	return s.drivers.Resolve(s.config.Driver)
}

// Gets the current driver name
func (s *Payment) GetDriverName() string {
	return string(s.config.Driver)
}

// CreatePaymentIntent creates a new order with a payment provider
func (s *Payment) CreatePaymentIntent(
	ctx context.Context,
	amount int64,
	currency string,
	customerID string,

	orderID string,
) (*PaymentIntentResult, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.CreatePaymentIntent(ctx, amount, currency, customerID, orderID)
}

// GetPaymentMethode
func (s *Payment) GetPaymentMethod(
	ctx context.Context,
	paymentMethodeID string,
) (string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return "", err
	}

	return driver.GetPaymentMethod(ctx, paymentMethodeID)
}

// UpdateCustomer
func (s *Payment) UpdateCustomer(
	ctx context.Context,
	customerID string,
	customer *Customer,
) error {
	driver, err := s.resolveDriver()

	if err != nil {
		return err
	}

	return driver.UpdateCustomer(ctx, customerID, customer)
}

// GetCustomerForPaymentIntent
func (s *Payment) GetCustomerForPaymentIntent(
	ctx context.Context,
	paymentIntentID string,
) (*Customer, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.GetCustomerForPaymentIntent(ctx, paymentIntentID)
}

// GetBillingDetailsForCharge
func (s *Payment) GetBillingDetailsForCharge(
	ctx context.Context,
	chargeID string,
) (*BillingDetails, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.GetBillingDetailsForCharge(ctx, chargeID)
}

// UpdatePaymentIntent updates an existing order with a payment provider
func (s *Payment) UpdatePaymentIntent(
	ctx context.Context,
	paymentIntentID string,
	amount int64,
	currency string,
	country string,
) (*PaymentIntentResult, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.UpdatePaymentIntent(ctx, paymentIntentID, amount, currency, country)
}

// CreateCustomer
func (s *Payment) CreateCustomer(
	ctx context.Context,
	customer *Customer,
) (string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return "", err
	}

	return driver.CreateCustomer(ctx, customer)
}

// CreateTransfer transfers funds from a given order to a vender
func (s *Payment) CreateTransfer(
	ctx context.Context,
	amount int64,
	currency string,
	recipientAccountID string,
	sourceTransactionID *string,
	transferGroupID string,
) (string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return "", err
	}

	return driver.CreateTransfer(
		ctx,
		amount,
		currency,
		recipientAccountID,
		sourceTransactionID,
		transferGroupID,
	)
}

// CreatePayout creates a new payout with a payment provider
func (s *Payment) CreatePayout(
	ctx context.Context,
	id string,
	amount int64,
	currency string,
	vendorAccountID string,
) (*string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.CreatePayout(ctx, id, amount, currency, vendorAccountID)
}

// CreateOrder creates a new order with a payment provider
func (s *Payment) CreateTransferReversal(
	ctx context.Context,
	amount int64,
	transferID string,
) (*string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.CreateTransferReversal(ctx, amount, transferID)
}

// CreateAccount creates a new account with a payment provider
func (s *Payment) CreateAccount(
	ctx context.Context,
	email string,
	// businessType string,
	payoutDelay int64,
	account Account,
) (string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return "", err
	}

	return driver.CreateAccount(ctx, email, payoutDelay, account)
}

// UpdateAccount updates an existing account with a payment provider
func (s *Payment) UpdateAccount(
	ctx context.Context,
	accountID string,
	account Account,
) error {
	driver, err := s.resolveDriver()

	if err != nil {
		return err
	}

	return driver.UpdateAccount(
		ctx,
		accountID,
		account,
	)
}

func (s *Payment) GetTransferForCharge(
	ctx context.Context,
	chargeID string,
) (*BalanceTransaction, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.GetBalanceTransactionForCharge(ctx, chargeID)
}

func (s *Payment) GetAccountVerificationLink(
	ctx context.Context,
	vendorAccountID string,
	returnUrl string,
	linkType string,
) (string, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return "", err
	}

	return driver.GetAccountVerificationLink(ctx, vendorAccountID, returnUrl, linkType)
}

// GetAccountWallet get the balance from the remote paymentAccount and creates the wallet domain object
func (s *Payment) GetAccountWallet(
	ctx context.Context,
	vendorAccountID string,
) (*Wallet, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.GetAccountWallet(ctx, vendorAccountID)
}

func (s *Payment) VerifyWebhookEvent(
	ctx context.Context,
	eventData string,
	signature string,
	webhookType WebhookType,
) (any, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.VerifyWebhookEvent(ctx, eventData, signature, webhookType)
}
