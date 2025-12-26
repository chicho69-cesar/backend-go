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
	GetAll() ([]*models.Loan, error)
	GetByID(id int64) (*models.Loan, error)
	GetByCode(code string) (*models.Loan, error)
	GetLoansFiltered(filter LoanFilter) ([]*models.Loan, error)
	Create(loan *models.Loan) (*models.Loan, error)
	Update(id int64, loan *models.Loan) (*models.Loan, error)
	Delete(id int64) error
}

type IReservationStore interface {
	GetAll() ([]*models.Reservation, error)
	GetByID(id int64) (*models.Reservation, error)
	GetActiveByUserAndBook(userID, bookID int64) (*models.Reservation, error)
	GetReservationsFiltered(filter ReservationFilter) ([]*models.Reservation, error)
	Create(reservation *models.Reservation) (*models.Reservation, error)
	Update(id int64, reservation *models.Reservation) (*models.Reservation, error)
	Delete(id int64) error
}

type IFineStore interface {
	GetAll() ([]*models.Fine, error)
	GetByID(id int64) (*models.Fine, error)
	GetFinesFiltered(filter FineFilter) ([]*models.Fine, error)
	Create(fine *models.Fine) (*models.Fine, error)
	Update(id int64, fine *models.Fine) (*models.Fine, error)
	Delete(id int64) error
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

func (s *LoanStore) GetAll() ([]*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id
		FROM loans 
		ORDER BY loan_date DESC
	`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (s *LoanStore) GetByID(id int64) (*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id
		FROM loans 
		WHERE id = ?
	`

	loan := &models.Loan{}

	err := s.db.
		QueryRow(query, id).
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
		)

	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanStore) GetByCode(code string) (*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id
		FROM loans 
		WHERE loan_code = ?
	`

	loan := &models.Loan{}

	err := s.db.
		QueryRow(query, code).
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
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanStore) GetLoansFiltered(filter LoanFilter) ([]*models.Loan, error) {
	query := `
		SELECT
			id, loan_code, user_id, copy_id, loan_date, due_date, 
			return_date, status, loan_days, renewals, notes, librarian_id
		FROM loans
	`

	var conditions []string
	var args []any

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
		)

		if err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (s *LoanStore) Create(loan *models.Loan) (*models.Loan, error) {
	query := `
		INSERT INTO loans (loan_code, user_id, copy_id, loan_date, due_date, return_date, status, loan_days, renewals, notes, librarian_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		loan.LoanCode, loan.UserID, loan.CopyID, loan.LoanDate, loan.DueDate,
		loan.ReturnDate, loan.Status, loan.LoanDays, loan.Renewals,
		loan.Notes, loan.LibrarianID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	loan.ID = id
	return loan, nil
}

func (s *LoanStore) Update(id int64, loan *models.Loan) (*models.Loan, error) {
	query := `
		UPDATE loans 
		SET
			loan_code = ?, user_id = ?, copy_id = ?, loan_date = ?, due_date = ?,
			return_date = ?, status = ?, loan_days = ?, renewals = ?, 
			notes = ?, librarian_id = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		loan.LoanCode, loan.UserID, loan.CopyID, loan.LoanDate, loan.DueDate,
		loan.ReturnDate, loan.Status, loan.LoanDays, loan.Renewals,
		loan.Notes, loan.LibrarianID, id,
	)

	if err != nil {
		return nil, err
	}

	loan.ID = id
	return loan, nil
}

func (s *LoanStore) Delete(id int64) error {
	query := `DELETE FROM loans WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ReservationStore) GetAll() ([]*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date, expiration_date, 
			status, priority, notified
		FROM reservations 
		ORDER BY reservation_date DESC
	`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (s *ReservationStore) GetByID(id int64) (*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date, expiration_date, 
			status, priority, notified
		FROM reservations 
		WHERE id = ?
	`

	reservation := &models.Reservation{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpirationDate,
			&reservation.Status,
			&reservation.Priority,
			&reservation.Notified,
		)

	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationStore) GetActiveByUserAndBook(userID, bookID int64) (*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date, expiration_date, 
			status, priority, notified
		FROM reservations 
		WHERE user_id = ? AND book_id = ? AND status IN ('Pending', 'Active')
		LIMIT 1
	`

	reservation := &models.Reservation{}

	err := s.db.
		QueryRow(query, userID, bookID).
		Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.BookID,
			&reservation.ReservationDate,
			&reservation.ExpirationDate,
			&reservation.Status,
			&reservation.Priority,
			&reservation.Notified,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *ReservationStore) GetReservationsFiltered(filter ReservationFilter) ([]*models.Reservation, error) {
	query := `
		SELECT
			id, user_id, book_id, reservation_date,
			expiration_date, status, priority, notified
		FROM reservations
	`

	var conditions []string
	var args []any

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
		)

		if err != nil {
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

func (s *ReservationStore) Create(reservation *models.Reservation) (*models.Reservation, error) {
	query := `
		INSERT INTO reservations (user_id, book_id, reservation_date, expiration_date, status, priority, notified)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		reservation.UserID, reservation.BookID, reservation.ReservationDate,
		reservation.ExpirationDate, reservation.Status, reservation.Priority,
		reservation.Notified,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	reservation.ID = id
	return reservation, nil
}

func (s *ReservationStore) Update(id int64, reservation *models.Reservation) (*models.Reservation, error) {
	query := `
		UPDATE reservations 
		SET
			user_id = ?, book_id = ?, reservation_date = ?,
			expiration_date = ?, status = ?, priority = ?, notified = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		reservation.UserID, reservation.BookID, reservation.ReservationDate,
		reservation.ExpirationDate, reservation.Status, reservation.Priority,
		reservation.Notified, id,
	)

	if err != nil {
		return nil, err
	}

	reservation.ID = id
	return reservation, nil
}

func (s *ReservationStore) Delete(id int64) error {
	query := `DELETE FROM reservations WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *FineStore) GetAll() ([]*models.Fine, error) {
	query := `
		SELECT
			id, user_id, loan_id, reason, amount,
			generated_date, payment_date, status, notes
		FROM fines 
		ORDER BY generated_date DESC
	`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		fines = append(fines, fine)
	}

	return fines, nil
}

func (s *FineStore) GetByID(id int64) (*models.Fine, error) {
	query := `
		SELECT
			id, user_id, loan_id, reason, amount,
			generated_date, payment_date, status, notes
		FROM fines 
		WHERE id = ?
	`

	fine := &models.Fine{}

	err := s.db.
		QueryRow(query, id).
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
		)

	if err != nil {
		return nil, err
	}

	return fine, nil
}

func (s *FineStore) GetFinesFiltered(filter FineFilter) ([]*models.Fine, error) {
	query := `
		SELECT
			id, user_id, loan_id, reason, amount, 
			generated_date, payment_date, status, notes
		FROM fines
	`

	var conditions []string
	var args []any

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
		)

		if err != nil {
			return nil, err
		}

		fines = append(fines, fine)
	}

	return fines, nil
}

func (s *FineStore) Create(fine *models.Fine) (*models.Fine, error) {
	query := `
		INSERT INTO fines (user_id, loan_id, reason, amount, generated_date, payment_date, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		fine.UserID, fine.LoanID, fine.Reason, fine.Amount, fine.GeneratedDate,
		fine.PaymentDate, fine.Status, fine.Notes,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	fine.ID = id
	return fine, nil
}

func (s *FineStore) Update(id int64, fine *models.Fine) (*models.Fine, error) {
	query := `
		UPDATE fines 
		SET 
			user_id = ?, loan_id = ?, reason = ?, amount = ?, 
			generated_date = ?, payment_date = ?, status = ?, notes = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		fine.UserID, fine.LoanID, fine.Reason, fine.Amount, fine.GeneratedDate,
		fine.PaymentDate, fine.Status, fine.Notes, id,
	)

	if err != nil {
		return nil, err
	}

	fine.ID = id
	return fine, nil
}

func (s *FineStore) Delete(id int64) error {
	query := `DELETE FROM fines WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
