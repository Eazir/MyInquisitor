package vo

import "time"

type Period struct {
	Start time.Time
	End   time.Time
}

func NewPeriod(start, end time.Time) (Period, error) {
	if start.After(end) {
		return Period{}, ErrInvalidPeriod
	}
	return Period{Start: start, End: end}, nil
}

var ErrInvalidPeriod = NewPeriodError("start date must be before end date")

type PeriodError string

func (e PeriodError) Error() string { return string(e) }

func NewPeriodError(text string) PeriodError { return PeriodError(text) }

func (p Period) MonthsBetween() int {
	months := (p.End.Year()-p.Start.Year())*12 + int(p.End.Month()-p.Start.Month())
	if months < 0 {
		return 0
	}
	return months
}
