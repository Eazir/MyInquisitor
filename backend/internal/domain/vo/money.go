package vo

import "errors"

type Money struct {
	amount float64
}

func NewMoney(amount float64) (Money, error) {
	if amount < 0 {
		return Money{}, errors.New("amount cannot be negative")
	}
	return Money{amount: amount}, nil
}

func (m Money) Amount() float64 { return m.amount }

func (m Money) Add(other Money) Money {
	return Money{amount: m.amount + other.amount}
}

func (m Money) Subtract(other Money) (Money, error) {
	if other.amount > m.amount {
		return Money{}, errors.New("insufficient funds")
	}
	return Money{amount: m.amount - other.amount}, nil
}

func (m Money) IsZero() bool { return m.amount == 0 }
