package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

type BookFilter struct {
	ISBN       string
	ShelfID    *int64
	AuthorID   *int64
	CategoryID *int64
}

type IBookStore interface {
	GetAll() ([]*models.Book, error)
	GetByID(id int64) (*models.Book, error)
	GetByISBN(isbn string) (*models.Book, error)
	GetBooksFiltered(filter BookFilter) ([]*models.Book, error)
	Create(book *models.Book) (*models.Book, error)
	Update(id int64, book *models.Book) (*models.Book, error)
	Delete(id int64) error

	GetBookAuthors(bookID int64) ([]*models.Author, error)
	AddAuthorToBook(bookAuthor *models.BookAuthor) error
	RemoveAuthorFromBook(bookID, authorID int64) error
	UpdateAuthorPosition(bookID, authorID int64, position int) error

	GetBookCategories(bookID int64) ([]*models.Category, error)
	AddCategoryToBook(bookCategory *models.BookCategory) error
	RemoveCategoryFromBook(bookID, categoryID int64) error
}

type BookStore struct {
	db *sql.DB
}

func NewBookStore(db *sql.DB) IBookStore {
	return &BookStore{
		db: db,
	}
}

func (s *BookStore) GetAll() ([]*models.Book, error) {
	query := `
		SELECT
			id, isbn, title, subtitle, edition, language, 
			publication_year, pages, synopsis, publisher_id, 
			shelf_id, status, registration_date
		FROM books
		ORDER BY title
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book

	for rows.Next() {
		book := &models.Book{}

		err := rows.Scan(
			&book.ID,
			&book.ISBN,
			&book.Title,
			&book.Subtitle,
			&book.Edition,
			&book.Language,
			&book.PublicationYear,
			&book.Pages,
			&book.Synopsis,
			&book.PublisherID,
			&book.ShelfID,
			&book.Status,
			&book.RegistrationDate,
		)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (s *BookStore) GetByID(id int64) (*models.Book, error) {
	query := `
		SELECT
			id, isbn, title, subtitle, edition, language, 
			publication_year, pages, synopsis, publisher_id, 
			shelf_id, status, registration_date
		FROM books 
		WHERE id = ?
	`

	book := &models.Book{}

	err := s.db.
		QueryRow(query, id).
		Scan(
			&book.ID,
			&book.ISBN,
			&book.Title,
			&book.Subtitle,
			&book.Edition,
			&book.Language,
			&book.PublicationYear,
			&book.Pages,
			&book.Synopsis,
			&book.PublisherID,
			&book.ShelfID,
			&book.Status,
			&book.RegistrationDate,
		)

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookStore) GetByISBN(isbn string) (*models.Book, error) {
	query := `
		SELECT
			id, isbn, title, subtitle, edition, language, 
			publication_year, pages, synopsis, publisher_id, 
			shelf_id, status, registration_date
		FROM books 
		WHERE isbn = ?
	`

	book := &models.Book{}

	err := s.db.
		QueryRow(query, isbn).
		Scan(
			&book.ID,
			&book.ISBN,
			&book.Title,
			&book.Subtitle,
			&book.Edition,
			&book.Language,
			&book.PublicationYear,
			&book.Pages,
			&book.Synopsis,
			&book.PublisherID,
			&book.ShelfID,
			&book.Status,
			&book.RegistrationDate,
		)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookStore) GetBooksFiltered(filter BookFilter) ([]*models.Book, error) {
	query := `
		SELECT DISTINCT
			b.id, b.isbn, b.title, b.subtitle, b.edition, b.language, 
			b.publication_year, b.pages, b.synopsis, b.publisher_id, 
			b.shelf_id, b.status, b.registration_date
		FROM books b
	`

	var joins []string

	if filter.AuthorID != nil {
		joins = append(joins, "INNER JOIN book_authors ba ON b.id = ba.book_id")
	}

	if filter.CategoryID != nil {
		joins = append(joins, "INNER JOIN book_categories bc ON b.id = bc.book_id")
	}

	if len(joins) > 0 {
		query += "\n" + strings.Join(joins, "\n")
	}

	var conditions []string
	var args []any

	if filter.ISBN != "" {
		conditions = append(conditions, "b.isbn = ?")
		args = append(args, filter.ISBN)
	}

	if filter.ShelfID != nil {
		conditions = append(conditions, "b.shelf_id = ?")
		args = append(args, *filter.ShelfID)
	}

	if filter.AuthorID != nil {
		conditions = append(conditions, "ba.author_id = ?")
		args = append(args, *filter.AuthorID)
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, "bc.category_id = ?")
		args = append(args, *filter.CategoryID)
	}

	if len(conditions) > 0 {
		query += "\nWHERE " + strings.Join(conditions, " AND ")
	}

	query += "\nORDER BY b.title"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book

	for rows.Next() {
		book := &models.Book{}

		err := rows.Scan(
			&book.ID,
			&book.ISBN,
			&book.Title,
			&book.Subtitle,
			&book.Edition,
			&book.Language,
			&book.PublicationYear,
			&book.Pages,
			&book.Synopsis,
			&book.PublisherID,
			&book.ShelfID,
			&book.Status,
			&book.RegistrationDate,
		)

		if err != nil {
			return nil, err
		}
		
		books = append(books, book)
	}

	return books, nil
}

func (s *BookStore) Create(book *models.Book) (*models.Book, error) {
	query := `
		INSERT INTO books (
			isbn, title, subtitle, edition, language, 
			publication_year, pages, synopsis, publisher_id, 
			shelf_id, status, registration_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		book.ISBN,
		book.Title,
		book.Subtitle,
		book.Edition,
		book.Language,
		book.PublicationYear,
		book.Pages,
		book.Synopsis,
		book.PublisherID,
		book.ShelfID,
		book.Status,
		book.RegistrationDate,
	)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	book.ID = id
	return book, nil
}

func (s *BookStore) Update(id int64, book *models.Book) (*models.Book, error) {
	query := `
		UPDATE books 
		SET
			isbn = ?, title = ?, subtitle = ?, edition = ?, language = ?,
			publication_year = ?, pages = ?, synopsis = ?, publisher_id = ?,
			shelf_id = ?, status = ?
		WHERE id = ?
	`

	_, err := s.db.Exec(
		query,
		book.ISBN,
		book.Title,
		book.Subtitle,
		book.Edition,
		book.Language,
		book.PublicationYear,
		book.Pages,
		book.Synopsis,
		book.PublisherID,
		book.ShelfID,
		book.Status,
		id,
	)

	if err != nil {
		return nil, err
	}

	book.ID = id
	return book, nil
}

func (s *BookStore) Delete(id int64) error {
	_, err := s.db.Exec("DELETE FROM book_authors WHERE book_id = ?", id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("DELETE FROM book_categories WHERE book_id = ?", id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookStore) GetBookAuthors(bookID int64) ([]*models.Author, error) {
	query := `
		SELECT a.id, a.first_name, a.last_name, a.biography, a.nationality
		FROM authors a
		INNER JOIN book_authors ba ON a.id = ba.author_id
		WHERE ba.book_id = ?
		ORDER BY ba.position
	`

	rows, err := s.db.Query(query, bookID)
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

func (s *BookStore) AddAuthorToBook(bookAuthor *models.BookAuthor) error {
	query := `SELECT COUNT(*) FROM book_authors WHERE book_id = ? AND author_id = ?`

	var count int

	err := s.db.QueryRow(
		query,
		bookAuthor.BookID,
		bookAuthor.AuthorID,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("El autor ya está asociado a este libro")
	}

	query = `INSERT INTO book_authors (book_id, author_id, position) VALUES (?, ?, ?)`

	_, err = s.db.Exec(query, bookAuthor.BookID, bookAuthor.AuthorID, bookAuthor.Position)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookStore) RemoveAuthorFromBook(bookID, authorID int64) error {
	query := `DELETE FROM book_authors WHERE book_id = ? AND author_id = ?`

	result, err := s.db.Exec(query, bookID, authorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("La relación libro-autor no existe")
	}

	return nil
}

func (s *BookStore) UpdateAuthorPosition(bookID, authorID int64, position int) error {
	query := `UPDATE book_authors SET position = ? WHERE book_id = ? AND author_id = ?`

	result, err := s.db.Exec(query, position, bookID, authorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("La relación libro-autor no existe")
	}

	return nil
}

func (s *BookStore) GetBookCategories(bookID int64) ([]*models.Category, error) {
	query := `
		SELECT c.id, c.name, c.description
		FROM categories c
		INNER JOIN book_categories bc ON c.id = bc.category_id
		WHERE bc.book_id = ?
		ORDER BY c.name
	`

	rows, err := s.db.Query(query, bookID)
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

func (s *BookStore) AddCategoryToBook(bookCategory *models.BookCategory) error {
	query := `SELECT COUNT(*) FROM book_categories WHERE book_id = ? AND category_id = ?`

	var count int

	err := s.db.QueryRow(
		query,
		bookCategory.BookID,
		bookCategory.CategoryID,
	).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("La categoría ya está asociada a este libro")
	}

	query = `INSERT INTO book_categories (book_id, category_id) VALUES (?, ?)`

	_, err = s.db.Exec(query, bookCategory.BookID, bookCategory.CategoryID)
	if err != nil {
		return err
	}

	return nil
}

func (s *BookStore) RemoveCategoryFromBook(bookID, categoryID int64) error {
	query := `DELETE FROM book_categories WHERE book_id = ? AND category_id = ?`

	result, err := s.db.Exec(query, bookID, categoryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("La relación libro-categoría no existe")
	}

	return nil
}
