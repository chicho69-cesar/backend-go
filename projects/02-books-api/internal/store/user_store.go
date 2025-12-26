package store

import (
	"database/sql"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type UserFilter struct {
	Code     string
	DNI      string
	UserType string
	Status   string
}

type IUserStore interface {
	GetAll(libraryID int64) ([]*models.User, error)
	GetByID(libraryID, id int64) (*models.User, error)
	GetByCode(libraryID int64, code string) (*models.User, error)
	GetByDNI(libraryID int64, dni string) (*models.User, error)
	GetUsersFiltered(libraryID int64, filter UserFilter) ([]*models.User, error)
	Create(libraryID int64, user *models.User) (*models.User, error)
	Update(libraryID, id int64, user *models.User) (*models.User, error)
	Delete(libraryID, id int64) error
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) IUserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetAll(libraryID int64) ([]*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date, library_id
		FROM users 
		WHERE library_id = ?
		ORDER BY last_name, first_name
	`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Code,
			&user.DNI,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.UserType,
			&user.Status,
			&user.RegistrationDate,
			&user.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) GetByID(libraryID, id int64) (*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date, library_id
		FROM users 
		WHERE id = ? AND library_id = ?
	`

	user := &models.User{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&user.ID,
			&user.Code,
			&user.DNI,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.UserType,
			&user.Status,
			&user.RegistrationDate,
			&user.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByCode(libraryID int64, code string) (*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date, library_id
		FROM users 
		WHERE code = ? AND library_id = ?
	`

	user := &models.User{}

	err := s.db.
		QueryRow(query, code, libraryID).
		Scan(
			&user.ID, &user.Code, &user.DNI, &user.FirstName, &user.LastName,
			&user.Email, &user.Phone, &user.Address, &user.UserType,
			&user.Status, &user.RegistrationDate, &user.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByDNI(libraryID int64, dni string) (*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date, library_id
		FROM users 
		WHERE dni = ? AND library_id = ?
	`

	user := &models.User{}

	err := s.db.
		QueryRow(query, dni, libraryID).
		Scan(
			&user.ID,
			&user.Code,
			&user.DNI,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.UserType,
			&user.Status,
			&user.RegistrationDate,
			&user.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetUsersFiltered(libraryID int64, filter UserFilter) ([]*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date, library_id
		FROM users
	`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.Code != "" {
		conditions = append(conditions, "code = ?")
		args = append(args, filter.Code)
	}

	if filter.DNI != "" {
		conditions = append(conditions, "dni = ?")
		args = append(args, filter.DNI)
	}

	if filter.UserType != "" {
		conditions = append(conditions, "user_type = ?")
		args = append(args, filter.UserType)
	}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY last_name, first_name"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Code,
			&user.DNI,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.Address,
			&user.UserType,
			&user.Status,
			&user.RegistrationDate,
			&user.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) Create(libraryID int64, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (code, dni, first_name, last_name, email, phone, address, user_type, status, registration_date, library_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		user.Code, user.DNI, user.FirstName, user.LastName, user.Email,
		user.Phone, user.Address, user.UserType, user.Status, user.RegistrationDate, libraryID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id
	user.LibraryID = libraryID

	return user, nil
}

func (s *UserStore) Update(libraryID, id int64, user *models.User) (*models.User, error) {
	query := `
		UPDATE users 
		SET
			code = ?, dni = ?, first_name = ?, last_name = ?, email = ?, 
			phone = ?, address = ?, user_type = ?, status = ?
		WHERE id = ? AND library_id = ?
	`

	_, err := s.db.Exec(
		query,
		user.Code, user.DNI, user.FirstName, user.LastName, user.Email,
		user.Phone, user.Address, user.UserType, user.Status, id, libraryID,
	)

	if err != nil {
		return nil, err
	}

	user.ID = id
	user.LibraryID = libraryID

	return user, nil
}

func (s *UserStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM users WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}
