package models

type Configuration struct {
	ID              int64   `json:"id"`
	StudentLoanDays int     `json:"student_loan_days"`
	TeacherLoanDays int     `json:"teacher_loan_days"`
	MaxRenewals     int     `json:"max_renewals"`
	MaxBooksPerLoan int     `json:"max_books_per_loan"`
	FinePerDay      float64 `json:"fine_per_day"`
	ReservationDays int     `json:"reservation_days"`
	GraceDays       int     `json:"grace_days"`
}
