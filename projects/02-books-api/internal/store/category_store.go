package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type ICategoryStore interface {
	GetAll(libraryID int64) ([]*models.Category, error)
	GetByID(libraryID, id int64) (*models.Category, error)
	Create(libraryID int64, category *models.Category) (*models.Category, error)
	Update(libraryID, id int64, category *models.Category) (*models.Category, error)
	Delete(libraryID, id int64) error
}

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) ICategoryStore {
	return &CategoryStore{
		db: db,
	}
}

func (s *CategoryStore) GetAll(libraryID int64) ([]*models.Category, error) {
	query := `SELECT id, name, description, library_id FROM categories WHERE library_id = ?`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		category := &models.Category{}

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (s *CategoryStore) GetByID(libraryID, id int64) (*models.Category, error) {
	query := `SELECT id, name, description, library_id FROM categories WHERE id = ? AND library_id = ?`

	var category = &models.Category{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryStore) Create(libraryID int64, category *models.Category) (*models.Category, error) {
	query := `INSERT INTO categories (name, description, library_id) VALUES (?, ?, ?)`

	result, err := s.db.Exec(query, category.Name, category.Description, libraryID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	category.ID = id
	category.LibraryID = libraryID

	return category, nil
}

func (s *CategoryStore) Update(libraryID, id int64, category *models.Category) (*models.Category, error) {
	query := `UPDATE categories SET name = ?, description = ? WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, category.Name, category.Description, id, libraryID)
	if err != nil {
		return nil, err
	}

	category.ID = id
	category.LibraryID = libraryID
	
	return category, nil
}

func (s *CategoryStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM categories WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}
