package store

import (
	"database/sql"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type LoanFilter struct {
	Code    string
	UserID  *int64
	CopyID  *int64
	Status  string
	Overdue bool
}

type ReservationFilter struct {
	UserID  *int64
	BookID  *int64
	Status  string
	Expired bool
}

type FineFilter struct {
	UserID  *int64
	LoanID  *int64
	Status  string
	Pending bool
}

type ILoanStore interface {
	GetAll(libraryID int64) ([]*models.Loan, error)
	GetByID(libraryID, id int64) (*models.Loan, error)
	GetByCode(libraryID int64, code string) (*models.Loan, error)
	GetLoansFiltered(libraryID int64, filter LoanFilter) ([]*models.Loan, error)
	Create(libraryID int64, loan *models.Loan) (*models.Loan, error)
	Update(libraryID, id int64, loan *models.Loan) (*models.Loan, error)
	Delete(libraryID, id int64) error
}

type IReservationStore interface {
	GetAll(libraryID int64) ([]*models.Reservation, error)
	GetByID(libraryID, id int64) (*models.Reservation, error)
	GetActiveByUserAndBook(libraryID, userID, bookID int64) (*models.Reservation, error)
	GetReservationsFiltered(libraryID int64, filter ReservationFilter) ([]*models.Reservation, error)
	Create(libraryID int64, reservation *models.Reservation) (*models.Reservation, error)
	Update(libraryID, id int64, reservation *models.Reservation) (*models.Reservation, error)
	Delete(libraryID, id int64) error
}

type IFineStore interface {
	GetAll(libraryID int64) ([]*models.Fine, error)
	GetByID(libraryID, id int64) (*models.Fine, error)
	GetFinesFiltered(libraryID int64, filter FineFilter) ([]*models.Fine, error)
	Create(libraryID int64, fine *models.Fine) (*models.Fine, error)
	Update(libraryID, id int64, fine *models.Fine) (*models.Fine, error)
	Delete(libraryID, id int64) error
}

type LoanStore struct {
	db *sql.DB
}

type ReservationStore struct {
	db *sql.DB
}

type FineStore struct {
	db *sql.DB
}

func NewLoanStore(db *sql.DB) ILoanStore {
	return &LoanStore{db: db}
}

func NewReservationStore(db *sql.DB) IReservationStore {
	return &ReservationStore{db: db}
}

func NewFineStore(db *sql.DB) IFineStore {
	return &FineStore{db: db}
}

func (s *LoanStore) GetAll(libraryID int64) ([]*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id, library_id 
		FROM loans 
		WHERE library_id = ? 
		ORDER BY loan_date DESC
	`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*models.Loan

	for rows.Next() {
		loan := &models.Loan{}

		err := rows.Scan(
			&loan.ID,
			&loan.LoanCode,
			&loan.UserID,
			&loan.CopyID,
			&loan.LoanDate,
			&loan.DueDate,
			&loan.ReturnDate,
			&loan.Status,
			&loan.LoanDays,
			&loan.Renewals,
			&loan.Notes,
			&loan.LibrarianID,
			&loan.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (s *LoanStore) GetByID(libraryID, id int64) (*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id, library_id 
		FROM loans 
		WHERE id = ? AND library_id = ?
	`

	loan := &models.Loan{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&loan.ID,
			&loan.LoanCode,
			&loan.UserID,
			&loan.CopyID,
			&loan.LoanDate,
			&loan.DueDate,
			&loan.ReturnDate,
			&loan.Status,
			&loan.LoanDays,
			&loan.Renewals,
			&loan.Notes,
			&loan.LibrarianID,
			&loan.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanStore) GetByCode(libraryID int64, code string) (*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id, library_id 
		FROM loans 
		WHERE loan_code = ? AND library_id = ?
	`

	loan := &models.Loan{}

	err := s.db.
		QueryRow(query, code, libraryID).
		Scan(
			&loan.ID,
			&loan.LoanCode,
			&loan.UserID,
			&loan.CopyID,
			&loan.LoanDate,
			&loan.DueDate,
			&loan.ReturnDate,
			&loan.Status,
			&loan.LoanDays,
			&loan.Renewals,
			&loan.Notes,
			&loan.LibrarianID,
			&loan.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanStore) GetLoansFiltered(libraryID int64, filter LoanFilter) ([]*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id, library_id 
		FROM loans
	`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.Code != "" {
		conditions = append(conditions, "loan_code = ?")
		args = append(args, filter.Code)
	}

	if filter.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.UserID)
	}

	if filter.CopyID != nil {
		conditions = append(conditions, "copy_id = ?")
		args = append(args, *filter.CopyID)
	}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Overdue {
		conditions = append(conditions, "status = 'Active' AND due_date < ?")
		args = append(args, time.Now())
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY loan_date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []*models.Loan

	for rows.Next() {
		loan := &models.Loan{}

		err := rows.Scan(
			&loan.ID,
			&loan.LoanCode,
			&loan.UserID,
			&loan.CopyID,
			&loan.LoanDate,
			&loan.DueDate,
			&loan.ReturnDate,
			&loan.Status,
			&loan.LoanDays,
			&loan.Renewals,
			&loan.Notes,
			&loan.LibrarianID,
			&loan.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (s *LoanStore) Create(libraryID int64, loan *models.Loan) (*models.Loan, error) {
	query := `
		INSERT INTO loans (loan_code, user_id, copy_id, loan_date, due_date, return_date, status, loan_days, renewals, notes, librarian_id, library_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		loan.LoanCode, loan.UserID, loan.CopyID, loan.LoanDate, loan.DueDate,
		loan.ReturnDate, loan.Status, loan.LoanDays, loan.Renewals,
		loan.Notes, loan.LibrarianID, libraryID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	loan.ID = id
	loan.LibraryID = libraryID

	return loan, nil
}

func (s *LoanStore) Update(libraryID, id int64, loan *models.Loan) (*models.Loan, error) {
	query := `
		UPDATE loans 
		SET
			loan_code = ?, user_id = ?, copy_id = ?, loan_date = ?, due_date = ?,
			return_date = ?, status = ?, loan_days = ?, renewals = ?, 
			notes = ?, librarian_id = ?
		WHERE id = ? AND library_id = ?
	`

	_, err := s.db.Exec(
		query,
		loan.LoanCode, loan.UserID, loan.CopyID, loan.LoanDate, loan.DueDate,
		loan.ReturnDate, loan.Status, loan.LoanDays, loan.Renewals,
		loan.Notes, loan.LibrarianID, id, libraryID,
	)

	if err != nil {
		return nil, err
	}

	loan.ID = id
	loan.LibraryID = libraryID

	return loan, nil
}

func (s *LoanStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM loans WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReservationStore) GetAll(libraryID int64) ([]*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date, expiration_date, 
			status, priority, notified, library_id
		FROM reservations 
		WHERE library_id = ? 
		ORDER BY reservation_date DESC
	`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*models.Reservation

	for rows.Next() {
		reservation := &models.Reservation{}

		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpirationDate,
			&reservation.Status,
			&reservation.Priority,
			&reservation.Notified,
			&reservation.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (s *ReservationStore) GetByID(libraryID, id int64) (*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date, expiration_date, 
			status, priority, notified, library_id 
		FROM reservations 
		WHERE id = ? AND library_id = ?
	`

	reservation := &models.Reservation{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpirationDate,
			&reservation.Status,
			&reservation.Priority,
			&reservation.Notified,
			&reservation.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationStore) GetActiveByUserAndBook(libraryID, userID, bookID int64) (*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date, expiration_date, 
			status, priority, notified, library_id
		FROM reservations 
		WHERE user_id = ? AND book_id = ? AND status IN ('Pending', 'Active') AND library_id = ?
		LIMIT 1
	`

	reservation := &models.Reservation{}

	err := s.db.
		QueryRow(query, userID, bookID, libraryID).
		Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpirationDate,
			&reservation.Status,
			&reservation.Priority,
			&reservation.Notified,
			&reservation.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationStore) GetReservationsFiltered(libraryID int64, filter ReservationFilter) ([]*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date,
			expiration_date, status, priority, notified, library_id
		FROM reservations 
	`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.UserID)
	}

	if filter.BookID != nil {
		conditions = append(conditions, "book_id = ?")
		args = append(args, *filter.BookID)
	}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Expired {
		conditions = append(conditions, "status IN ('Pending', 'Active') AND expiration_date < ?")
		args = append(args, time.Now())
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY priority DESC, reservation_date ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*models.Reservation

	for rows.Next() {
		reservation := &models.Reservation{}

		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpirationDate,
			&reservation.Status,
			&reservation.Priority,
			&reservation.Notified,
			&reservation.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (s *ReservationStore) Create(libraryID int64, reservation *models.Reservation) (*models.Reservation, error) {
	query := `
		INSERT INTO reservations (user_id, book_id, reservation_date, expiration_date, status, priority, notified, library_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		reservation.UserID, reservation.BookID, reservation.ReservationDate,
		reservation.ExpirationDate, reservation.Status, reservation.Priority,
		reservation.Notified, libraryID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	reservation.ID = id
	reservation.LibraryID = libraryID

	return reservation, nil
}

func (s *ReservationStore) Update(libraryID, id int64, reservation *models.Reservation) (*models.Reservation, error) {
	query := `
		UPDATE reservations 
		SET
			user_id = ?, book_id = ?, reservation_date = ?,
			expiration_date = ?, status = ?, priority = ?, notified = ?
		WHERE id = ? AND library_id = ?
	`

	_, err := s.db.Exec(
		query,
		reservation.UserID, reservation.BookID, reservation.ReservationDate,
		reservation.ExpirationDate, reservation.Status, reservation.Priority,
		reservation.Notified, id, libraryID,
	)

	if err != nil {
		return nil, err
	}

	reservation.ID = id
	reservation.LibraryID = libraryID

	return reservation, nil
}

func (s *ReservationStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM reservations WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *FineStore) GetAll(libraryID int64) ([]*models.Fine, error) {
	query := `
		SELECT
			id, user_id, loan_id, reason, amount,
			generated_date, payment_date, status, notes, library_id
		FROM fines 
		WHERE library_id = ? 
		ORDER BY generated_date DESC
	`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fines []*models.Fine

	for rows.Next() {
		fine := &models.Fine{}

		err := rows.Scan(
			&fine.ID,
			&fine.UserID,
			&fine.LoanID,
			&fine.Reason,
			&fine.Amount,
			&fine.GeneratedDate,
			&fine.PaymentDate,
			&fine.Status,
			&fine.Notes,
			&fine.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		fines = append(fines, fine)
	}

	return fines, nil
}

func (s *FineStore) GetByID(libraryID, id int64) (*models.Fine, error) {
	query := `
		SELECT
			id, user_id, loan_id, reason, amount,
			generated_date, payment_date, status, notes, library_id
		FROM fines 
		WHERE id = ? AND library_id = ?
	`

	fine := &models.Fine{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&fine.ID,
			&fine.UserID,
			&fine.LoanID,
			&fine.Reason,
			&fine.Amount,
			&fine.GeneratedDate,
			&fine.PaymentDate,
			&fine.Status,
			&fine.Notes,
			&fine.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return fine, nil
}

func (s *FineStore) GetFinesFiltered(libraryID int64, filter FineFilter) ([]*models.Fine, error) {
	query := `
		SELECT
			id, user_id, loan_id, reason, amount, 
			generated_date, payment_date, status, notes, library_id
		FROM fines
	`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.UserID != nil {
		conditions = append(conditions, "user_id = ?")
		args = append(args, *filter.UserID)
	}

	if filter.LoanID != nil {
		conditions = append(conditions, "loan_id = ?")
		args = append(args, *filter.LoanID)
	}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Pending {
		conditions = append(conditions, "status = 'Pending'")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY generated_date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fines []*models.Fine

	for rows.Next() {
		fine := &models.Fine{}

		err := rows.Scan(
			&fine.ID,
			&fine.UserID,
			&fine.LoanID,
			&fine.Reason,
			&fine.Amount,
			&fine.GeneratedDate,
			&fine.PaymentDate,
			&fine.Status,
			&fine.Notes,
			&fine.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		fines = append(fines, fine)
	}

	return fines, nil
}

func (s *FineStore) Create(libraryID int64, fine *models.Fine) (*models.Fine, error) {
	query := `
		INSERT INTO fines (user_id, loan_id, reason, amount, generated_date, payment_date, status, notes, library_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		fine.UserID, fine.LoanID, fine.Reason, fine.Amount, fine.GeneratedDate,
		fine.PaymentDate, fine.Status, fine.Notes, libraryID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	fine.ID = id
	fine.LibraryID = libraryID

	return fine, nil
}

func (s *FineStore) Update(libraryID, id int64, fine *models.Fine) (*models.Fine, error) {
	query := `
		UPDATE fines 
		SET 
			user_id = ?, loan_id = ?, reason = ?, amount = ?, 
			generated_date = ?, payment_date = ?, status = ?, notes = ?
		WHERE id = ? AND library_id = ?
	`

	_, err := s.db.Exec(
		query,
		fine.UserID, fine.LoanID, fine.Reason, fine.Amount, fine.GeneratedDate,
		fine.PaymentDate, fine.Status, fine.Notes, id, libraryID,
	)

	if err != nil {
		return nil, err
	}

	fine.ID = id
	fine.LibraryID = libraryID

	return fine, nil
}

func (s *FineStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM fines WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}
