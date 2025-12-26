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
	GetAll() ([]*models.User, error)
	GetByID(id int64) (*models.User, error)
	GetByCode(code string) (*models.User, error)
	GetByDNI(dni string) (*models.User, error)
	GetUsersFiltered(filter UserFilter) ([]*models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(id int64, user *models.User) (*models.User, error)
	Delete(id int64) error
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) IUserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetAll() ([]*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date
		FROM users 
		ORDER BY last_name, first_name
	`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date
		FROM users 
		WHERE id = ?
	`

	user := &models.User{}

	err := s.db.
		QueryRow(query, id).
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
		)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByCode(code string) (*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date
		FROM users 
		WHERE code = ?
	`

	user := &models.User{}

	err := s.db.
		QueryRow(query, code).
		Scan(
			&user.ID, &user.Code, &user.DNI, &user.FirstName, &user.LastName,
			&user.Email, &user.Phone, &user.Address, &user.UserType,
			&user.Status, &user.RegistrationDate,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByDNI(dni string) (*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date
		FROM users 
		WHERE dni = ?
	`

	user := &models.User{}

	err := s.db.
		QueryRow(query, dni).
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
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetUsersFiltered(filter UserFilter) ([]*models.User, error) {
	query := `
		SELECT
			id, code, dni, first_name, last_name, email, phone, 
			address, user_type, status, registration_date
		FROM users
	`

	var conditions []string
	var args []any

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
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *UserStore) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (code, dni, first_name, last_name, email, phone, address, user_type, status, registration_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		user.Code, user.DNI, user.FirstName, user.LastName, user.Email,
		user.Phone, user.Address, user.UserType, user.Status, user.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserStore) Update(id int64, user *models.User) (*models.User, error) {
	query := `
		UPDATE users 
		SET
			code = ?, dni = ?, first_name = ?, last_name = ?, email = ?, 
			phone = ?, address = ?, user_type = ?, status = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		user.Code, user.DNI, user.FirstName, user.LastName, user.Email,
		user.Phone, user.Address, user.UserType, user.Status, id,
	)

	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (s *UserStore) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
