package payment

type Balance struct {
	amount   int64
	currency string
}

type Wallet struct {
	Available Balance
	Pending   Balance
}
