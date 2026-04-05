package models

import "time"

type RecordType string

const (
	RecordTypeIncome  RecordType = "income"
	RecordTypeExpense RecordType = "expense"
)

type FinancialRecord struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Amount    string     `json:"amount"`
	Type      RecordType `json:"type"`
	Category  string     `json:"category"`
	Date      time.Time  `json:"date"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Summary struct {
	TotalIncome   string `json:"total_income"`
	TotalExpenses string `json:"total_expenses"`
	NetBalance    string `json:"net_balance"`
}

type CategoryTotal struct {
	Category string `json:"category"`
	Total    string `json:"total"`
}

type TrendPoint struct {
	Period  string `json:"period"`
	Income  string `json:"income"`
	Expense string `json:"expense"`
}
