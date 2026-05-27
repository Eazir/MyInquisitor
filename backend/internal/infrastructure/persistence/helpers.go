package persistence

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toPGUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func toPGText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func toPGDate(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

func toPGDatePtr(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func toPGNumeric(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	if err := n.Scan(f); err != nil {
		return pgtype.Numeric{Valid: false}
	}
	return n
}

func fromPGNumeric(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	f, err := n.Float64Value()
	if err != nil {
		return 0
	}
	val := f.Float64
	if math.IsNaN(val) {
		return 0
	}
	return val
}

func fromPGText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func fromPGUUID(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	id := uuid.UUID(u.Bytes)
	return &id
}

func fromPGDate(d pgtype.Date) time.Time {
	if !d.Valid {
		return time.Time{}
	}
	return d.Time
}

func fromPGDatePtr(d pgtype.Date) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}

func fromPGTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
