package payment

import (
	"github.com/shopspring/decimal"
)

type BusinessType string

const (
	BusinessTypeIndividual BusinessType = "individual"
	BusinessTypeCompany    BusinessType = "company"
	BusinessTypeNonProfit  BusinessType = "non_profit"
)

type Country string

const (
	CountryUS Country = "US"
	CountryGB Country = "GB"
	CountryDE Country = "DE"
	CountryFR Country = "FR"
	CountryBE Country = "BE"
	CountryNL Country = "NL"
	CountryES Country = "ES"
	CountryIE Country = "IE"
)

type WebhookType string

const (
	WebhookTypeConnect WebhookType = "connect"
	WebhookTypeAccount WebhookType = "account"
)

type PaymentIntentResult struct {
	ClientSecret    string
	PaymentIntentID string
}

type BalanceTransaction struct {
	Amount       decimal.Decimal
	Currency     string
	ProviderFee  decimal.Decimal
	ExchangeRate decimal.Decimal
}

type Account struct {
	Email              *string
	AddressCity        *string
	AddressCountryCode *string
	AddressStreet      *string
	AddressZip         *string
	BusinessName       *string
	BusinessVATID      *string
	CountryCode        *string
	BusinessType       BusinessType
}

type Customer struct {
	ID    string
	Email string
}

type BillingDetails struct {
	ID          string
	CountryCode string
	Name        string
	Street      string
	Zip         string
	City        string
	Email       string
}
