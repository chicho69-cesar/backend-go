package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type ICategoryStore interface {
	GetAll() ([]*models.Category, error)
	GetByID(id int64) (*models.Category, error)
	Create(category *models.Category) (*models.Category, error)
	Update(id int64, category *models.Category) (*models.Category, error)
	Delete(id int64) error
}

type CategoryStore struct {
	db *sql.DB
}

func NewCategoryStore(db *sql.DB) ICategoryStore {
	return &CategoryStore{
		db: db,
	}
}

func (s *CategoryStore) GetAll() ([]*models.Category, error) {
	query := `SELECT id, name, description FROM categories`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (s *CategoryStore) GetByID(id int64) (*models.Category, error) {
	query := `SELECT id, name, description FROM categories WHERE id = ?`

	var category = &models.Category{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&category.ID,
			&category.Name,
			&category.Description,
		)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryStore) Create(category *models.Category) (*models.Category, error) {
	query := `INSERT INTO categories (name, description) VALUES (?, ?)`

	result, err := s.db.Exec(query, category.Name, category.Description)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	category.ID = id

	return category, nil
}

func (s *CategoryStore) Update(id int64, category *models.Category) (*models.Category, error) {
	query := `UPDATE categories SET name = ?, description = ? WHERE id = ?`

	_, err := s.db.Exec(query, category.Name, category.Description, id)
	if err != nil {
		return nil, err
	}

	category.ID = id

	return category, nil
}

func (s *CategoryStore) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
