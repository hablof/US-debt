package model

import "time"

type DebtDifference struct {
	Valid bool
	Grows bool
	Value uint64
}

type Debt struct {
	Amuont uint64
	Date   time.Time

	Diff DebtDifference
}
