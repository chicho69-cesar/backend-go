package models

import (
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
)

type Fine struct {
	ID            int64               `json:"id"`
	UserID        int64               `json:"user_id"`
	LoanID        database.NullInt64  `json:"loan_id"`
	Reason        string              `json:"reason"` // Overdue, Damage, Loss
	Amount        float64             `json:"amount"`
	GeneratedDate time.Time           `json:"generated_date"`
	PaymentDate   database.NullTime   `json:"payment_date"`
	Status        string              `json:"status"` // Pending, Paid, Waived
	Notes         database.NullString `json:"notes"`
}
