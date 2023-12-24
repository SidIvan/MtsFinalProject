package model

import "github.com/rmg/iso4217"

type Money struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
}

func (m *Money) IsValid() bool {
	return m.IsAmountValid() && m.IsCurrencyValid()
}

func (m *Money) IsAmountValid() bool {
	return IsAmountValid(m.Amount)
}

func (m *Money) IsCurrencyValid() bool {
	return IsCurrencyValid(m.Currency)
}

func IsAmountValid(amount float64) bool {
	return amount > 0
}

func IsCurrencyValid(currency string) bool {
	code, _ := iso4217.ByName(currency)
	return code > 0
}
