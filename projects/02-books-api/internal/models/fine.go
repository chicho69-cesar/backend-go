package models

import (
	"database/sql"
	"time"
)

type Fine struct {
	ID            int64          `json:"id"`
	UserID        int64          `json:"user_id"`
	LoanID        sql.NullInt64  `json:"loan_id"`
	Reason        string         `json:"reason"` // Overdue, Damage, Loss
	Amount        float64        `json:"amount"`
	GeneratedDate time.Time      `json:"generated_date"`
	PaymentDate   sql.NullTime   `json:"payment_date"`
	Status        string         `json:"status"` // Pending, Paid, Waived
	Notes         sql.NullString `json:"notes"`
}
