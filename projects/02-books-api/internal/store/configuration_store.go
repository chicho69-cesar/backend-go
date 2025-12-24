package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type ConfigStore interface {
	GetCurrent() (*models.Configuration, error)
	Update(config *models.Configuration) (*models.Configuration, error)
}

type configurationStore struct {
	db *sql.DB
}

func NewConfigurationStore(db *sql.DB) ConfigStore {
	return &configurationStore{
		db: db,
	}
}

func (s *configurationStore) GetCurrent() (*models.Configuration, error) {
	query := `SELECT id, student_loan_days, teacher_loan_days, max_renewals, max_books_per_loan, fine_per_day, reservation_days, grace_days FROM configuration LIMIT 1`

	config := &models.Configuration{}

	err := s.db.
		QueryRow(query).
		Scan(
			&config.ID,
			&config.StudentLoanDays,
			&config.TeacherLoanDays,
			&config.MaxRenewals,
			&config.MaxBooksPerLoan,
			&config.FinePerDay,
			&config.ReservationDays,
			&config.GraceDays,
		)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *configurationStore) Update(config *models.Configuration) (*models.Configuration, error) {
	current, err := s.GetCurrent()
	if err != nil {
		return nil, err
	}

	query := `UPDATE configuration SET 
		student_loan_days = ?, 
		teacher_loan_days = ?, 
		max_renewals = ?, 
		max_books_per_loan = ?, 
		fine_per_day = ?, 
		reservation_days = ?, 
		grace_days = ? 
		WHERE id = ?`

	_, err = s.db.Exec(
		query,
		config.StudentLoanDays,
		config.TeacherLoanDays,
		config.MaxRenewals,
		config.MaxBooksPerLoan,
		config.FinePerDay,
		config.ReservationDays,
		config.GraceDays,
		current.ID,
	)

	if err != nil {
		return nil, err
	}

	config.ID = current.ID

	return config, nil
}
