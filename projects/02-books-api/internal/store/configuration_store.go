package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type IConfigStore interface {
	GetByLibraryID(libraryID int64) (*models.Configuration, error)
	Update(libraryID int64, config *models.Configuration) (*models.Configuration, error)
	Create(libraryID int64, config *models.Configuration) (*models.Configuration, error)
}

type ConfigurationStore struct {
	db *sql.DB
}

func NewConfigurationStore(db *sql.DB) IConfigStore {
	return &ConfigurationStore{
		db: db,
	}
}

func (s *ConfigurationStore) GetByLibraryID(libraryID int64) (*models.Configuration, error) {
	query := `SELECT id, student_loan_days, teacher_loan_days, max_renewals, max_books_per_loan, fine_per_day, reservation_days, grace_days, library_id FROM configuration WHERE library_id = ?`

	config := &models.Configuration{}

	err := s.db.
		QueryRow(query, libraryID).
		Scan(
			&config.ID,
			&config.StudentLoanDays,
			&config.TeacherLoanDays,
			&config.MaxRenewals,
			&config.MaxBooksPerLoan,
			&config.FinePerDay,
			&config.ReservationDays,
			&config.GraceDays,
			&config.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *ConfigurationStore) Update(libraryID int64, config *models.Configuration) (*models.Configuration, error) {
	query := `
		UPDATE configuration 
		SET 
			student_loan_days = ?, teacher_loan_days = ?, max_renewals = ?, 
			max_books_per_loan = ?, fine_per_day = ?, reservation_days = ?, grace_days = ? 
		WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(
		query,
		config.StudentLoanDays,
		config.TeacherLoanDays,
		config.MaxRenewals,
		config.MaxBooksPerLoan,
		config.FinePerDay,
		config.ReservationDays,
		config.GraceDays,
		config.ID,
		libraryID,
	)

	if err != nil {
		return nil, err
	}

	config.LibraryID = libraryID

	return config, nil
}

func (s *ConfigurationStore) Create(libraryID int64, config *models.Configuration) (*models.Configuration, error) {
	query := `
		INSERT INTO configuration 
		(student_loan_days, teacher_loan_days, max_renewals, max_books_per_loan, fine_per_day, reservation_days, grace_days, library_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.Exec(
		query,
		config.StudentLoanDays,
		config.TeacherLoanDays,
		config.MaxRenewals,
		config.MaxBooksPerLoan,
		config.FinePerDay,
		config.ReservationDays,
		config.GraceDays,
		libraryID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	config.ID = id
	config.LibraryID = libraryID
	
	return config, nil
}
