package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type IPublisherStore interface {
	GetAll() ([]*models.Publisher, error)
	GetByID(id int64) (*models.Publisher, error)
	Create(publisher *models.Publisher) (*models.Publisher, error)
	Update(id int64, publisher *models.Publisher) (*models.Publisher, error)
	Delete(id int64) error
}

type PublisherStore struct {
	db *sql.DB
}

func NewPublisherStore(db *sql.DB) IPublisherStore {
	return &PublisherStore{
		db: db,
	}
}

func (s *PublisherStore) GetAll() ([]*models.Publisher, error) {
	query := `SELECT id, name, country FROM publishers`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publishers []*models.Publisher

	for rows.Next() {
		publisher := &models.Publisher{}

		err := rows.Scan(
			&publisher.ID,
			&publisher.Name,
			&publisher.Country,
		)

		if err != nil {
			return nil, err
		}

		publishers = append(publishers, publisher)
	}

	return publishers, nil
}

func (s *PublisherStore) GetByID(id int64) (*models.Publisher, error) {
	query := `SELECT id, name, country FROM publishers WHERE id = ?`

	var publisher = &models.Publisher{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&publisher.ID,
			&publisher.Name,
			&publisher.Country,
		)

	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func (s *PublisherStore) Create(publisher *models.Publisher) (*models.Publisher, error) {
	query := `INSERT INTO publishers (name, country) VALUES (?, ?)`

	result, err := s.db.Exec(query, publisher.Name, publisher.Country)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	publisher.ID = id

	return publisher, nil
}

func (s *PublisherStore) Update(id int64, publisher *models.Publisher) (*models.Publisher, error) {
	query := `UPDATE publishers SET name = ?, country = ? WHERE id = ?`

	_, err := s.db.Exec(query, publisher.Name, publisher.Country, id)
	if err != nil {
		return nil, err
	}

	publisher.ID = id

	return publisher, nil
}

func (s *PublisherStore) Delete(id int64) error {
	query := `DELETE FROM publishers WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
