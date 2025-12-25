package store

import (
	"database/sql"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type IAuthorStore interface {
	GetAll() ([]*models.Author, error)
	GetByID(id int64) (*models.Author, error)
	Create(author *models.Author) (*models.Author, error)
	Update(id int64, author *models.Author) (*models.Author, error)
	Delete(id int64) error
}

type AuthorStore struct {
	db *sql.DB
}

func NewAuthorStore(db *sql.DB) IAuthorStore {
	return &AuthorStore{
		db: db,
	}
}

func (s *AuthorStore) GetAll() ([]*models.Author, error) {
	query := `SELECT id, first_name, last_name, biography, nationality FROM authors`

	rows, err := s.db.Query(query)
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
		)

		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func (s *AuthorStore) GetByID(id int64) (*models.Author, error) {
	query := `SELECT id, first_name, last_name, biography, nationality FROM authors WHERE id = ?`

	author := &models.Author{}

	err := s.db.
		QueryRow(query, id).
		Scan(&author.ID, &author.FirstName, &author.LastName, &author.Biography, &author.Nationality)

	if err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorStore) Create(author *models.Author) (*models.Author, error) {
	query := `INSERT INTO authors (first_name, last_name, biography, nationality) VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, author.FirstName, author.LastName, author.Biography, author.Nationality)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	author.ID = id

	return author, nil
}

func (s *AuthorStore) Update(id int64, author *models.Author) (*models.Author, error) {
	query := `UPDATE authors SET first_name = ?, last_name = ?, biography = ?, nationality = ? WHERE id = ?`

	_, err := s.db.Exec(query, author.FirstName, author.LastName, author.Biography, author.Nationality, id)
	if err != nil {
		return nil, err
	}

	author.ID = id

	return author, nil
}

func (s *AuthorStore) Delete(id int64) error {
	query := `DELETE FROM authors WHERE id = ?`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
