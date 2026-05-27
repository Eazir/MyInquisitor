package vo

type DebtStatus string

const (
	DebtStatusActive  DebtStatus = "active"
	DebtStatusPaid    DebtStatus = "paid"
	DebtStatusSettled DebtStatus = "settled"
)

func (s DebtStatus) String() string { return string(s) }

func NewDebtStatus(status string) (DebtStatus, error) {
	s := DebtStatus(status)
	switch s {
	case DebtStatusActive, DebtStatusPaid, DebtStatusSettled:
		return s, nil
	default:
		return "", ErrInvalidDebtStatus
	}
}

var ErrInvalidDebtStatus = NewDebtStatusError("invalid debt status: must be active, paid, or settled")

type DebtStatusError string

func (e DebtStatusError) Error() string { return string(e) }

func NewDebtStatusError(text string) DebtStatusError { return DebtStatusError(text) }
