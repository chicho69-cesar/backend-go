package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type IPublisherStore interface {
	GetAll(libraryID int64) ([]*models.Publisher, error)
	GetByID(libraryID, id int64) (*models.Publisher, error)
	Create(libraryID int64, publisher *models.Publisher) (*models.Publisher, error)
	Update(libraryID, id int64, publisher *models.Publisher) (*models.Publisher, error)
	Delete(libraryID, id int64) error
}

type PublisherStore struct {
	db *sql.DB
}

func NewPublisherStore(db *sql.DB) IPublisherStore {
	return &PublisherStore{
		db: db,
	}
}

func (s *PublisherStore) GetAll(libraryID int64) ([]*models.Publisher, error) {
	query := `SELECT id, name, country, library_id FROM publishers WHERE library_id = ?`

	rows, err := s.db.Query(query, libraryID)
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
			&publisher.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		publishers = append(publishers, publisher)
	}

	return publishers, nil
}

func (s *PublisherStore) GetByID(libraryID, id int64) (*models.Publisher, error) {
	query := `SELECT id, name, country, library_id FROM publishers WHERE id = ? AND library_id = ?`

	var publisher = &models.Publisher{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&publisher.ID,
			&publisher.Name,
			&publisher.Country,
			&publisher.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func (s *PublisherStore) Create(libraryID int64, publisher *models.Publisher) (*models.Publisher, error) {
	query := `INSERT INTO publishers (name, country, library_id) VALUES (?, ?, ?)`

	result, err := s.db.Exec(query, publisher.Name, publisher.Country, libraryID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	publisher.ID = id
	publisher.LibraryID = libraryID

	return publisher, nil
}

func (s *PublisherStore) Update(libraryID, id int64, publisher *models.Publisher) (*models.Publisher, error) {
	query := `UPDATE publishers SET name = ?, country = ? WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, publisher.Name, publisher.Country, id, libraryID)
	if err != nil {
		return nil, err
	}

	publisher.ID = id
	publisher.LibraryID = libraryID

	return publisher, nil
}

func (s *PublisherStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM publishers WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}
