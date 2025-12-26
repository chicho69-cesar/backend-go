package store

import (
	"database/sql"
	"strings"

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

type ILibraryZoneStore interface {
	GetAll() ([]*models.LibraryZone, error)
	GetByID(id int64) (*models.LibraryZone, error)
	GetByCode(code string) (*models.LibraryZone, error)
	GetZonesFiltered(filter LibraryZoneFilter) ([]*models.LibraryZone, error)
	Create(zone *models.LibraryZone) (*models.LibraryZone, error)
	Update(id int64, zone *models.LibraryZone) (*models.LibraryZone, error)
	Delete(id int64) error
}

type IShelfStore interface {
	GetAll() ([]*models.Shelf, error)
	GetByID(id int64) (*models.Shelf, error)
	GetByCode(code string) (*models.Shelf, error)
	GetShelvesFiltered(filter ShelfFilter) ([]*models.Shelf, error)
	Create(shelf *models.Shelf) (*models.Shelf, error)
	Update(id int64, shelf *models.Shelf) (*models.Shelf, error)
	Delete(id int64) error
}

type ICopyStore interface {
	GetAll() ([]*models.Copy, error)
	GetByID(id int64) (*models.Copy, error)
	GetByCode(code string) (*models.Copy, error)
	GetCopiesFiltered(filter CopyFilter) ([]*models.Copy, error)
	Create(copy *models.Copy) (*models.Copy, error)
	Update(id int64, copy *models.Copy) (*models.Copy, error)
	Delete(id int64) error
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

func NewLibraryZoneStore(db *sql.DB) ILibraryZoneStore {
	return &LibraryZoneStore{db: db}
}

func NewShelfStore(db *sql.DB) IShelfStore {
	return &ShelfStore{db: db}
}

func NewCopyStore(db *sql.DB) ICopyStore {
	return &CopyStore{db: db}
}

func (s *LibraryZoneStore) GetAll() ([]*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor FROM library_zones ORDER BY code`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		zones = append(zones, zone)
	}

	return zones, nil
}

func (s *LibraryZoneStore) GetByID(id int64) (*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor FROM library_zones WHERE id = ?`

	zone := &models.LibraryZone{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&zone.ID,
			&zone.Code,
			&zone.Name,
			&zone.Description,
			&zone.Floor,
		)

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *LibraryZoneStore) GetByCode(code string) (*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor FROM library_zones WHERE code = ?`

	zone := &models.LibraryZone{}

	err := s.db.
		QueryRow(query, code).
		Scan(
			&zone.ID,
			&zone.Code,
			&zone.Name,
			&zone.Description,
			&zone.Floor,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *LibraryZoneStore) GetZonesFiltered(filter LibraryZoneFilter) ([]*models.LibraryZone, error) {
	query := `SELECT id, code, name, description, floor FROM library_zones`

	var conditions []string
	var args []any

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
		)
		if err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}

	return zones, nil
}

func (s *LibraryZoneStore) Create(zone *models.LibraryZone) (*models.LibraryZone, error) {
	query := `INSERT INTO library_zones (code, name, description, floor) VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, zone.Code, zone.Name, zone.Description, zone.Floor)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	zone.ID = id
	return zone, nil
}

func (s *LibraryZoneStore) Update(id int64, zone *models.LibraryZone) (*models.LibraryZone, error) {
	query := `UPDATE library_zones SET code = ?, name = ?, description = ?, floor = ? WHERE id = ?`

	_, err := s.db.Exec(query, zone.Code, zone.Name, zone.Description, zone.Floor, id)
	if err != nil {
		return nil, err
	}

	zone.ID = id
	return zone, nil
}

func (s *LibraryZoneStore) Delete(id int64) error {
	query := `DELETE FROM library_zones WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ShelfStore) GetAll() ([]*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description FROM shelves ORDER BY code`

	rows, err := s.db.Query(query)
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
			)

		if err != nil {
			return nil, err
		}

		shelves = append(shelves, shelf)
	}

	return shelves, nil
}

func (s *ShelfStore) GetByID(id int64) (*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description FROM shelves WHERE id = ?`

	shelf := &models.Shelf{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&shelf.ID,
			&shelf.Code,
			&shelf.ZoneID,
			&shelf.Description,
		)

	if err != nil {
		return nil, err
	}

	return shelf, nil
}

func (s *ShelfStore) GetByCode(code string) (*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description FROM shelves WHERE code = ?`

	shelf := &models.Shelf{}

	err := s.db.
		QueryRow(query, code).
		Scan(
			&shelf.ID,
			&shelf.Code,
			&shelf.ZoneID,
			&shelf.Description,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return shelf, nil
}

func (s *ShelfStore) GetShelvesFiltered(filter ShelfFilter) ([]*models.Shelf, error) {
	query := `SELECT id, code, zone_id, description FROM shelves`

	var conditions []string
	var args []any

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
		)

		if err != nil {
			return nil, err
		}

		shelves = append(shelves, shelf)
	}

	return shelves, nil
}

func (s *ShelfStore) Create(shelf *models.Shelf) (*models.Shelf, error) {
	query := `INSERT INTO shelves (code, zone_id, description) VALUES (?, ?, ?)`

	result, err := s.db.Exec(query, shelf.Code, shelf.ZoneID, shelf.Description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	shelf.ID = id
	return shelf, nil
}

func (s *ShelfStore) Update(id int64, shelf *models.Shelf) (*models.Shelf, error) {
	query := `UPDATE shelves SET code = ?, zone_id = ?, description = ? WHERE id = ?`

	_, err := s.db.Exec(query, shelf.Code, shelf.ZoneID, shelf.Description, id)
	if err != nil {
		return nil, err
	}

	shelf.ID = id
	return shelf, nil
}

func (s *ShelfStore) Delete(id int64) error {
	query := `DELETE FROM shelves WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CopyStore) GetAll() ([]*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes
		FROM copies 
		ORDER BY code
	`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		copies = append(copies, copy)
	}

	return copies, nil
}

func (s *CopyStore) GetByID(id int64) (*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes
		FROM copies 
		WHERE id = ?
	`

	copy := &models.Copy{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&copy.ID,
			&copy.Code,
			&copy.BookID,
			&copy.Status,
			&copy.Condition,
			&copy.AcquisitionDate,
			&copy.PurchasePrice,
			&copy.Notes,
		)

	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (s *CopyStore) GetByCode(code string) (*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes
		FROM copies 
		WHERE code = ?
	`

	copy := &models.Copy{}

	err := s.db.
		QueryRow(query, code).
		Scan(
			&copy.ID,
			&copy.Code,
			&copy.BookID,
			&copy.Status,
			&copy.Condition,
			&copy.AcquisitionDate,
			&copy.PurchasePrice,
			&copy.Notes,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return copy, nil
}

func (s *CopyStore) GetCopiesFiltered(filter CopyFilter) ([]*models.Copy, error) {
	query := `
		SELECT
			id, code, book_id, status, condition, 
			acquisition_date, purchase_price, notes
		FROM copies
	`

	var conditions []string
	var args []any

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
		)

		if err != nil {
			return nil, err
		}
		
		copies = append(copies, copy)
	}

	return copies, nil
}

func (s *CopyStore) Create(copy *models.Copy) (*models.Copy, error) {
	query := `
		INSERT INTO copies (code, book_id, status, condition, acquisition_date, purchase_price, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?)
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
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	copy.ID = id
	return copy, nil
}

func (s *CopyStore) Update(id int64, copy *models.Copy) (*models.Copy, error) {
	query := `
		UPDATE copies 
		SET
			code = ?, book_id = ?, status = ?, condition = ?, 
			acquisition_date = ?, purchase_price = ?, notes = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		copy.Code,
		copy.BookID,
		copy.Status,
		copy.Condition,
		copy.AcquisitionDate,
		copy.PurchasePrice,
		copy.Notes, id,
	)

	if err != nil {
		return nil, err
	}

	copy.ID = id
	return copy, nil
}

func (s *CopyStore) Delete(id int64) error {
	query := `DELETE FROM copies WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
