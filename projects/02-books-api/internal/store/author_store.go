package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type IAuthorStore interface {
	GetAll(libraryID int64) ([]*models.Author, error)
	GetByID(libraryID, id int64) (*models.Author, error)
	Create(libraryID int64, author *models.Author) (*models.Author, error)
	Update(libraryID, id int64, author *models.Author) (*models.Author, error)
	Delete(libraryID, id int64) error
}

type AuthorStore struct {
	db *sql.DB
}

func NewAuthorStore(db *sql.DB) IAuthorStore {
	return &AuthorStore{
		db: db,
	}
}

func (s *AuthorStore) GetAll(libraryID int64) ([]*models.Author, error) {
	query := `SELECT id, first_name, last_name, biography, nationality, library_id FROM authors WHERE library_id = ?`

	rows, err := s.db.Query(query, libraryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*models.Author

	for rows.Next() {
		author := &models.Author{}

		err := rows.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Biography,
			&author.Nationality,
			&author.LibraryID,
		)

		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (s *AuthorStore) GetByID(libraryID, id int64) (*models.Author, error) {
	query := `SELECT id, first_name, last_name, biography, nationality, library_id FROM authors WHERE id = ? AND library_id = ?`

	author := &models.Author{}

	err := s.db.
		QueryRow(query, id, libraryID).
		Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Biography,
			&author.Nationality,
			&author.LibraryID,
		)

	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorStore) Create(libraryID int64, author *models.Author) (*models.Author, error) {
	query := `INSERT INTO authors (first_name, last_name, biography, nationality, library_id) VALUES (?, ?, ?, ?, ?)`

	result, err := s.db.Exec(query, author.FirstName, author.LastName, author.Biography, author.Nationality, libraryID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	author.ID = id
	author.LibraryID = libraryID

	return author, nil
}

func (s *AuthorStore) Update(libraryID, id int64, author *models.Author) (*models.Author, error) {
	query := `UPDATE authors SET first_name = ?, last_name = ?, biography = ?, nationality = ? WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, author.FirstName, author.LastName, author.Biography, author.Nationality, id, libraryID)
	if err != nil {
		return nil, err
	}

	author.ID = id
	author.LibraryID = libraryID
	
	return author, nil
}

func (s *AuthorStore) Delete(libraryID, id int64) error {
	query := `DELETE FROM authors WHERE id = ? AND library_id = ?`

	_, err := s.db.Exec(query, id, libraryID)
	if err != nil {
		return err
	}

	return nil
}
