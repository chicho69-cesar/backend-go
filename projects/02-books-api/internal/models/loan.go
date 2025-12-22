package models

import (
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
)

type Loan struct {
	ID          int64               `json:"id"`
	LoanCode    string              `json:"loan_code"`
	UserID      int64               `json:"user_id"`
	CopyID      int64               `json:"copy_id"`
	LoanDate    time.Time           `json:"loan_date"`
	DueDate     time.Time           `json:"due_date"`
	ReturnDate  database.NullTime   `json:"return_date"`
	Status      string              `json:"status"` // Active, Returned, Overdue, Lost
	LoanDays    int                 `json:"loan_days"`
	Renewals    int                 `json:"renewals"`
	Notes       database.NullString `json:"notes"`
	LibrarianID database.NullInt64  `json:"librarian_id"`
}
