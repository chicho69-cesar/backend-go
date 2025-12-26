package store

import (
	"database/sql"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type LibraryZoneFilter struct {
	Code  string
	Floor *int
}

type ShelfFilter struct {
	Code   string
	ZoneID *int64
}

type CopyFilter struct {
	Code      string
	BookID    *int64
	Status    string
	Condition string
}

type ILibraryStore interface {
	GetAll() ([]*models.Library, error)
	GetByID(id int64) (*models.Library, error)
	GetByUsername(username string) (*models.Library, error)
	CheckPassword(username, password string) (bool, error)
	Create(library *models.Library) (*models.Library, error)
	Update(id int64, library *models.Library) (*models.Library, error)
	Delete(id int64) error
}

type ILibraryZoneStore interface {
	GetAll(libraryID int64) ([]*models.LibraryZone, error)
	GetByID(libraryID, id int64) (*models.LibraryZone, error)
	GetByCode(libraryID int64, code string) (*models.LibraryZone, error)
	GetZonesFiltered(libraryID int64, filter LibraryZoneFilter) ([]*models.LibraryZone, error)
	Create(libraryID int64, zone *models.LibraryZone) (*models.LibraryZone, error)
	Update(libraryID, id int64, zone *models.LibraryZone) (*models.LibraryZone, error)
	Delete(libraryID, id int64) error
}

type IShelfStore interface {
	GetAll(libraryID int64) ([]*models.Shelf, error)
	GetByID(libraryID, id int64) (*models.Shelf, error)
	GetByCode(libraryID int64, code string) (*models.Shelf, error)
	GetShelvesFiltered(libraryID int64, filter ShelfFilter) ([]*models.Shelf, error)
	Create(libraryID int64, shelf *models.Shelf) (*models.Shelf, error)
	Update(libraryID, id int64, shelf *models.Shelf) (*models.Shelf, error)
	Delete(libraryID, id int64) error
}

type ICopyStore interface {
	GetAll(libraryID int64) ([]*models.Copy, error)
	GetByID(libraryID, id int64) (*models.Copy, error)
	GetByCode(libraryID int64, code string) (*models.Copy, error)
	GetCopiesFiltered(libraryID int64, filter CopyFilter) ([]*models.Copy, error)
	Create(libraryID int64, copy *models.Copy) (*models.Copy, error)
	Update(libraryID, id int64, copy *models.Copy) (*models.Copy, error)
	Delete(libraryID, id int64) error
}

type LibraryStore struct {
	db *sql.DB
}

type LibraryZoneStore struct {
	db *sql.DB
}

type ShelfStore struct {
	db *sql.DB
}

type CopyStore struct {
	db *sql.DB
}

func NewLibraryStore(db *sql.DB) ILibraryStore {
	return &LibraryStore{db: db}
}

func NewLibraryZoneStore(db *sql.DB) ILibraryZoneStore {
	return &LibraryZoneStore{db: db}
}

func NewShelfStore(db *sql.DB) IShelfStore {
	return &ShelfStore{db: db}
}

func NewCopyStore(db *sql.DB) ICopyStore {
	return &CopyStore{db: db}
}

func (s *LibraryStore) GetAll() ([]*models.Library, error) {
	query := `SELECT id, name, address, city, state, zip_code, country, phone, email, website FROM libraries ORDER BY name`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libraries []*models.Library

	for rows.Next() {
		library := &models.Library{}

		err := rows.Scan(
			&library.ID,
			&library.Name,
			&library.Address,
			&library.City,
			&library.State,
			&library.ZipCode,
			&library.Country,
			&library.Phone,
			&library.Email,
			&library.Website,
		)

		if err != nil {
			return nil, err
		}

		libraries = append(libraries, library)
	}

	return libraries, nil
}

func (s *LibraryStore) GetByID(id int64) (*models.Library, error) {
	query := `SELECT id, name, address, city, state, zip_code, country, phone, email, website FROM libraries WHERE id = ?`

	library := &models.Library{}
	err := s.db.
		QueryRow(query, id).
		Scan(
			&library.ID,
			&library.Name,
			&library.Address,
			&library.City,
			&library.State,
			&library.ZipCode,
			&library.Country,
			&library.Phone,
			&library.Email,
			&library.Website,
		)

	if err != nil {
		return nil, err
	}

	return library, nil
}

func (s *LibraryStore) GetByUsername(username string) (*models.Library, error) {
	query := `SELECT id, name, address, city, state, zip_code, country, phone, email, website, username FROM libraries WHERE username = ?`

	library := &models.Library{}
	err := s.db.
		QueryRow(query, username).
		Scan(
			&library.ID,
			&library.Name,
			&library.Address,
			&library.City,
			&library.State,
			&library.ZipCode,
			&library.Country,
			&library.Phone,
			&library.Email,
			&library.Website,
			&library.Username,
		)

	if err != nil {
		return nil, err
	}

	return library, nil
}

func (s *LibraryStore) CheckPassword(username, password string) (bool, error) {
	query := `SELECT password FROM libraries WHERE username = ?`

	var hashedPassword string

	err := s.db.
		QueryRow(query, username).
		Scan(&hashedPassword)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	isValid := database.CheckPasswordHash(password, hashedPassword)
	return isValid, nil
}

func (s *LibraryStore) Create(library *models.Library) (*models.Library, error) {
	query := `INSERT INTO libraries (name, address, city, state, zip_code, country, phone, email, website, username, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	newPassword, err := database.HashPassword(library.Password)
	if err != nil {
		return nil, err
	}

	library.Password = newPassword

	result, err := s.db.Exec(
		query,
		library.Name, library.Address, library.City, library.State,
		library.ZipCode, library.Country, library.Phone, library.Email,
		library.Website, library.Username, library.Password,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	library.ID = id
	library.Password = ""

	configQuery := `INSERT INTO configuration (student_loan_days, teacher_loan_days, max_renewals, max_books_per_loan, fine_per_day, reservation_days, grace_days, library_id) VALUES (15, 30, 2, 5, 0.50, 3, 2, ?)`

	_, err = s.db.Exec(configQuery, id)
	if err != nil {
		return nil, err
	}

	return library, nil
}

func (s *LibraryStore) Update(id int64, library *models.Library) (*models.Library, error) {
	query := `UPDATE libraries SET name = ?, address = ?, city = ?, state = ?, zip_code = ?, country = ?, phone = ?, email = ?, website = ? WHERE id = ?`

	_, err := s.db.Exec(
		query,
		library.Name, library.Address, library.City, library.State,
		library.ZipCode, library.Country, library.Phone, library.Email,
		library.Website, id,
	)
	if err != nil {
		return nil, err
	}

	library.ID = id
	return nil, nil
}

func (s *LibraryStore) Delete(id int64) error {
	query := `DELETE FROM libraries WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *LibraryZoneStore) GetAll(libraryID int64) ([]*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor, library_id FROM library_zones WHERE library_id = ? ORDER BY code`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []*models.LibraryZone

	for rows.Next() {
		zone := &models.LibraryZone{}

		err := rows.Scan(
			&zone.ID,
			&zone.Code,
			&zone.Name,
			&zone.Description,
			&zone.Floor,
			&zone.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

func (s *LibraryZoneStore) GetByID(libraryID, id int64) (*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor, library_id FROM library_zones WHERE id = ? AND library_id = ?`

	zone := &models.LibraryZone{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&zone.ID,
			&zone.Code,
			&zone.Name,
			&zone.Description,
			&zone.Floor,
			&zone.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *LibraryZoneStore) GetByCode(libraryID int64, code string) (*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor, library_id FROM library_zones WHERE code = ? AND library_id = ?`

	zone := &models.LibraryZone{}

	err := s.db.
		QueryRow(query, code, libraryID).
		Scan(
			&zone.ID,
			&zone.Code,
			&zone.Name,
			&zone.Description,
			&zone.Floor,
			&zone.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *LibraryZoneStore) GetZonesFiltered(libraryID int64, filter LibraryZoneFilter) ([]*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor, library_id FROM library_zones`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.Code != "" {
		conditions = append(conditions, "code = ?")
		args = append(args, filter.Code)
	}

	if filter.Floor != nil {
		conditions = append(conditions, "floor = ?")
		args = append(args, *filter.Floor)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY code"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []*models.LibraryZone

	for rows.Next() {
		zone := &models.LibraryZone{}

		err := rows.Scan(
			&zone.ID,
			&zone.Code,
			&zone.Name,
			&zone.Description,
			&zone.Floor,
			&zone.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

func (s *LibraryZoneStore) Create(libraryID int64, zone *models.LibraryZone) (*models.LibraryZone, error) {
	query := `INSERT INTO library_zones (code, name, description, floor, library_id) VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, zone.Code, zone.Name, zone.Description, zone.Floor, libraryID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	zone.ID = id
	zone.LibraryID = libraryID

	return zone, nil
}

func (s *LibraryZoneStore) Update(libraryID, id int64, zone *models.LibraryZone) (*models.LibraryZone, error) {
	query := `UPDATE library_zones SET code = ?, name = ?, description = ?, floor = ? WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, zone.Code, zone.Name, zone.Description, zone.Floor, id, libraryID)
	if err != nil {
		return nil, err
	}

	zone.ID = id
	zone.LibraryID = libraryID

	return zone, nil
}

func (s *LibraryZoneStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM library_zones WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *ShelfStore) GetAll(libraryID int64) ([]*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description, library_id FROM shelves WHERE library_id = ? ORDER BY code`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelves []*models.Shelf

	for rows.Next() {
		shelf := &models.Shelf{}

		err := rows.
			Scan(
				&shelf.ID,
				&shelf.Code,
				&shelf.ZoneID,
				&shelf.Description,
				&shelf.LibraryID,
			)

		if err != nil {
			return nil, err
		}

		shelves = append(shelves, shelf)
	}

	return shelves, nil
}

func (s *ShelfStore) GetByID(libraryID, id int64) (*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description, library_id FROM shelves WHERE id = ? AND library_id = ?`

	shelf := &models.Shelf{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&shelf.ID,
			&shelf.Code,
			&shelf.ZoneID,
			&shelf.Description,
			&shelf.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return shelf, nil
}

func (s *ShelfStore) GetByCode(libraryID int64, code string) (*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description, library_id FROM shelves WHERE code = ? AND library_id = ?`

	shelf := &models.Shelf{}

	err := s.db.
		QueryRow(query, code).
		Scan(
			&shelf.ID,
			&shelf.Code,
			&shelf.ZoneID,
			&shelf.Description,
			&shelf.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return shelf, nil
}

func (s *ShelfStore) GetShelvesFiltered(libraryID int64, filter ShelfFilter) ([]*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description, library_id FROM shelves`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.Code != "" {
		conditions = append(conditions, "code = ?")
		args = append(args, filter.Code)
	}

	if filter.ZoneID != nil {
		conditions = append(conditions, "zone_id = ?")
		args = append(args, *filter.ZoneID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY code"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shelves []*models.Shelf

	for rows.Next() {
		shelf := &models.Shelf{}

		err := rows.Scan(
			&shelf.ID,
			&shelf.Code,
			&shelf.ZoneID,
			&shelf.Description,
			&shelf.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		shelves = append(shelves, shelf)
	}

	return shelves, nil
}

func (s *ShelfStore) Create(libraryID int64, shelf *models.Shelf) (*models.Shelf, error) {
	query := `INSERT INTO shelves (code, zone_id, description, library_id) VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, shelf.Code, shelf.ZoneID, shelf.Description, libraryID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	shelf.ID = id
	shelf.LibraryID = libraryID

	return shelf, nil
}

func (s *ShelfStore) Update(libraryID, id int64, shelf *models.Shelf) (*models.Shelf, error) {
	query := `UPDATE shelves SET code = ?, zone_id = ?, description = ? WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, shelf.Code, shelf.ZoneID, shelf.Description, id, libraryID)
	if err != nil {
		return nil, err
	}

	shelf.ID = id
	shelf.LibraryID = libraryID

	return shelf, nil
}

func (s *ShelfStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM shelves WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *CopyStore) GetAll(libraryID int64) ([]*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes, library_id
		FROM copies 
		WHERE library_id = ? 
		ORDER BY code
	`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var copies []*models.Copy

	for rows.Next() {
		copy := &models.Copy{}

		err := rows.Scan(
			&copy.ID,
			&copy.Code,
			&copy.BookID,
			&copy.Status,
			&copy.Condition,
			&copy.AcquisitionDate,
			&copy.PurchasePrice,
			&copy.Notes,
			&copy.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		copies = append(copies, copy)
	}

	return copies, nil
}

func (s *CopyStore) GetByID(libraryID, id int64) (*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes, library_id 
		FROM copies 
		WHERE id = ? AND library_id = ?
	`

	copy := &models.Copy{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&copy.ID,
			&copy.Code,
			&copy.BookID,
			&copy.Status,
			&copy.Condition,
			&copy.AcquisitionDate,
			&copy.PurchasePrice,
			&copy.Notes,
			&copy.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (s *CopyStore) GetByCode(libraryID int64, code string) (*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes, library_id 
		FROM copies 
		WHERE code = ? AND library_id = ?
	`

	copy := &models.Copy{}

	err := s.db.
		QueryRow(query, code, libraryID).
		Scan(
			&copy.ID,
			&copy.Code,
			&copy.BookID,
			&copy.Status,
			&copy.Condition,
			&copy.AcquisitionDate,
			&copy.PurchasePrice,
			&copy.Notes,
			&copy.LibraryID,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (s *CopyStore) GetCopiesFiltered(libraryID int64, filter CopyFilter) ([]*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes, library_id
		FROM copies
	`

	var conditions []string
	var args []any

	conditions = append(conditions, "library_id = ?")
	args = append(args, libraryID)

	if filter.Code != "" {
		conditions = append(conditions, "code = ?")
		args = append(args, filter.Code)
	}

	if filter.BookID != nil {
		conditions = append(conditions, "book_id = ?")
		args = append(args, *filter.BookID)
	}

	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}

	if filter.Condition != "" {
		conditions = append(conditions, "condition = ?")
		args = append(args, filter.Condition)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY code"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var copies []*models.Copy

	for rows.Next() {
		copy := &models.Copy{}

		err := rows.Scan(
			&copy.ID,
			&copy.Code,
			&copy.BookID,
			&copy.Status,
			&copy.Condition,
			&copy.AcquisitionDate,
			&copy.PurchasePrice,
			&copy.Notes,
			&copy.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		copies = append(copies, copy)
	}

	return copies, nil
}

func (s *CopyStore) Create(libraryID int64, copy *models.Copy) (*models.Copy, error) {
	query := `
		INSERT INTO copies (code, book_id, status, condition, acquisition_date, purchase_price, notes, library_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		copy.Code,
		copy.BookID,
		copy.Status,
		copy.Condition,
		copy.AcquisitionDate,
		copy.PurchasePrice,
		copy.Notes,
		libraryID,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	copy.ID = id
	copy.LibraryID = libraryID

	return copy, nil
}

func (s *CopyStore) Update(libraryID, id int64, copy *models.Copy) (*models.Copy, error) {
	query := `
		UPDATE copies 
		SET
			code = ?, book_id = ?, status = ?, condition = ?, 
			acquisition_date = ?, purchase_price = ?, notes = ?
		WHERE id = ? AND library_id = ?
	`

	_, err := s.db.Exec(
		query,
		copy.Code,
		copy.BookID,
		copy.Status,
		copy.Condition,
		copy.AcquisitionDate,
		copy.PurchasePrice,
		copy.Notes,
		id,
		libraryID,
	)

	if err != nil {
		return nil, err
	}

	copy.ID = id
	copy.LibraryID = libraryID

	return copy, nil
}

func (s *CopyStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM copies WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}
