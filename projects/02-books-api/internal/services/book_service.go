package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/validations"
)

type BookService struct {
	bookStore        store.IBookStore
	authorStore      store.IAuthorStore
	copyStore        store.ICopyStore
	reservationStore store.IReservationStore
}

func NewBookService(bookStore store.IBookStore, authorStore store.IAuthorStore, copyStore store.ICopyStore, reservationStore store.IReservationStore) *BookService {
	return &BookService{
		bookStore:        bookStore,
		authorStore:      authorStore,
		copyStore:        copyStore,
		reservationStore: reservationStore,
	}
}

func (s *BookService) GetAllBooks(libraryID int64) ([]*models.Book, error) {
	books, err := s.bookStore.GetAll(libraryID)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los libros: %w", err)
	}

	return books, nil
}

func (s *BookService) GetBookByID(libraryID, id int64) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("El ID del libro es inválido")
	}

	book, err := s.bookStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el libro con ID %d: %w", id, err)
	}

	return book, nil
}

func (s *BookService) GetBookByISBN(libraryID int64, isbn string) (*models.Book, error) {
	if strings.TrimSpace(isbn) == "" {
		return nil, errors.New("El ISBN no puede estar vacío")
	}

	book, err := s.bookStore.GetByISBN(libraryID, isbn)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el libro con ISBN %s: %w", isbn, err)
	}

	if book == nil {
		return nil, fmt.Errorf("No se encontró un libro con ISBN %s", isbn)
	}

	return book, nil
}

func (s *BookService) GetBooksFiltered(libraryID int64, filter store.BookFilter) ([]*models.Book, error) {
	if filter.ShelfID != nil {
		if *filter.ShelfID <= 0 {
			return nil, errors.New("El ID del estante es inválido")
		}
	}

	if filter.AuthorID != nil {
		if *filter.AuthorID <= 0 {
			return nil, errors.New("El ID del autor es inválido")
		}

		_, err := s.authorStore.GetByID(libraryID, *filter.AuthorID)
		if err != nil {
			return nil, fmt.Errorf("El autor con ID %d no existe: %w", *filter.AuthorID, err)
		}
	}

	if filter.CategoryID != nil {
		if *filter.CategoryID <= 0 {
			return nil, errors.New("El ID de la categoría es inválido")
		}
	}

	if filter.ISBN != "" {
		filter.ISBN = strings.TrimSpace(filter.ISBN)

		if filter.ISBN == "" {
			return nil, errors.New("El ISBN no puede estar vacío")
		}
	}

	books, err := s.bookStore.GetBooksFiltered(libraryID, filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los libros filtrados: %w", err)
	}

	return books, nil
}

func (s *BookService) CreateBook(libraryID int64, book *models.Book) (*models.Book, error) {
	if err := validations.ValidateBook(book); err != nil {
		return nil, fmt.Errorf("validación fallida: %w", err)
	}

	existingBook, _ := s.bookStore.GetByISBN(libraryID, book.ISBN)
	if existingBook != nil {
		return nil, fmt.Errorf("Ya existe un libro con el ISBN %s", book.ISBN)
	}

	book.LibraryID = libraryID
	book.ISBN = strings.TrimSpace(book.ISBN)
	book.Title = strings.TrimSpace(book.Title)

	if book.Subtitle.Valid {
		book.Subtitle.String = strings.TrimSpace(book.Subtitle.String)
	}

	if book.Edition.Valid {
		book.Edition.String = strings.TrimSpace(book.Edition.String)
	}

	if book.Language.Valid {
		book.Language.String = strings.TrimSpace(book.Language.String)
	}

	if book.Synopsis.Valid {
		book.Synopsis.String = strings.TrimSpace(book.Synopsis.String)
	}

	if book.RegistrationDate.IsZero() {
		book.RegistrationDate = time.Now()
	}

	if strings.TrimSpace(book.Status) == "" {
		book.Status = "Available"
	}

	createdBook, err := s.bookStore.Create(libraryID, book)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el libro: %w", err)
	}

	return createdBook, nil
}

func (s *BookService) UpdateBook(libraryID, id int64, book *models.Book) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("El ID del libro es inválido")
	}

	existingBook, err := s.bookStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("El libro con ID %d no existe: %w", id, err)
	}

	if existingBook == nil {
		return nil, fmt.Errorf("El libro con ID %d no fue encontrado", id)
	}

	if err := validations.ValidateBook(book); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	bookWithISBN, _ := s.bookStore.GetByISBN(libraryID, book.ISBN)
	if bookWithISBN != nil && bookWithISBN.ID != id {
		return nil, fmt.Errorf("Ya existe otro libro con el ISBN %s", book.ISBN)
	}

	book.ISBN = strings.TrimSpace(book.ISBN)
	book.Title = strings.TrimSpace(book.Title)

	if book.Subtitle.Valid {
		book.Subtitle.String = strings.TrimSpace(book.Subtitle.String)
	}

	if book.Edition.Valid {
		book.Edition.String = strings.TrimSpace(book.Edition.String)
	}

	if book.Language.Valid {
		book.Language.String = strings.TrimSpace(book.Language.String)
	}

	if book.Synopsis.Valid {
		book.Synopsis.String = strings.TrimSpace(book.Synopsis.String)
	}

	book.RegistrationDate = existingBook.RegistrationDate

	updatedBook, err := s.bookStore.Update(libraryID, id, book)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar el libro con ID %d: %w", id, err)
	}

	return updatedBook, nil
}

func (s *BookService) DeleteBook(libraryID, id int64) error {
	if id <= 0 {
		return errors.New("El ID del libro es inválido")
	}

	existingBook, err := s.bookStore.GetByID(libraryID, id)
	if err != nil {
		return fmt.Errorf("El libro con ID %d no existe: %w", id, err)
	}

	if existingBook == nil {
		return fmt.Errorf("El libro con ID %d no fue encontrado", id)
	}

	copies, err := s.copyStore.GetCopiesFiltered(libraryID, store.CopyFilter{
		BookID: &id,
		Status: "Borrowed",
	})
	if err == nil && len(copies) > 0 {
		return fmt.Errorf("No se puede eliminar el libro porque tiene %d copia(s) prestada(s)", len(copies))
	}

	activeReservations, err := s.reservationStore.GetReservationsFiltered(libraryID, store.ReservationFilter{
		BookID: &id,
		Status: "Pending",
	})
	if err == nil && len(activeReservations) > 0 {
		return fmt.Errorf("No se puede eliminar el libro porque tiene %d reservación(es) activa(s)", len(activeReservations))
	}

	processingReservations, err := s.reservationStore.GetReservationsFiltered(libraryID, store.ReservationFilter{
		BookID: &id,
		Status: "Active",
	})
	if err == nil && len(processingReservations) > 0 {
		return fmt.Errorf("No se puede eliminar el libro porque tiene %d reservación(es) en proceso", len(processingReservations))
	}

	if err := s.bookStore.Delete(libraryID, id); err != nil {
		return fmt.Errorf("Error al eliminar el libro con ID %d: %w", id, err)
	}

	return nil
}

func (s *BookService) GetBookAuthors(libraryID, bookID int64) ([]*models.Author, error) {
	if bookID <= 0 {
		return nil, errors.New("El ID del libro es inválido")
	}

	_, err := s.bookStore.GetByID(libraryID, bookID)
	if err != nil {
		return nil, fmt.Errorf("El libro con ID %d no existe: %w", bookID, err)
	}

	authors, err := s.bookStore.GetBookAuthors(libraryID, bookID)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los autores del libro: %w", err)
	}

	return authors, nil
}

func (s *BookService) AddAuthorToBook(libraryID int64, bookAuthor *models.BookAuthor) error {
	if err := validations.ValidateBookAuthor(bookAuthor); err != nil {
		return fmt.Errorf("Validación fallida: %w", err)
	}

	_, err := s.bookStore.GetByID(libraryID, bookAuthor.BookID)
	if err != nil {
		return fmt.Errorf("El libro con ID %d no existe: %w", bookAuthor.BookID, err)
	}

	_, err = s.authorStore.GetByID(libraryID, bookAuthor.AuthorID)
	if err != nil {
		return fmt.Errorf("El autor con ID %d no existe: %w", bookAuthor.AuthorID, err)
	}

	if err := s.bookStore.AddAuthorToBook(libraryID, bookAuthor); err != nil {
		return fmt.Errorf("Error al agregar el autor al libro: %w", err)
	}

	return nil
}

func (s *BookService) RemoveAuthorFromBook(libraryID, bookID, authorID int64) error {
	if bookID <= 0 {
		return errors.New("El ID del libro es inválido")
	}

	if authorID <= 0 {
		return errors.New("El ID del autor es inválido")
	}

	_, err := s.bookStore.GetByID(libraryID, bookID)
	if err != nil {
		return fmt.Errorf("El libro con ID %d no existe: %w", bookID, err)
	}

	if err := s.bookStore.RemoveAuthorFromBook(libraryID, bookID, authorID); err != nil {
		return fmt.Errorf("Error al eliminar el autor del libro: %w", err)
	}

	return nil
}

func (s *BookService) UpdateAuthorPosition(libraryID, bookID, authorID int64, position int) error {
	if bookID <= 0 {
		return errors.New("El ID del libro es inválido")
	}

	if authorID <= 0 {
		return errors.New("El ID del autor es inválido")
	}

	if position < 1 {
		return errors.New("La posición debe ser al menos 1")
	}

	_, err := s.bookStore.GetByID(libraryID, bookID)
	if err != nil {
		return fmt.Errorf("El libro con ID %d no existe: %w", bookID, err)
	}

	if err := s.bookStore.UpdateAuthorPosition(libraryID, bookID, authorID, position); err != nil {
		return fmt.Errorf("Error al actualizar la posición del autor: %w", err)
	}

	return nil
}

func (s *BookService) GetBookCategories(libraryID, bookID int64) ([]*models.Category, error) {
	if bookID <= 0 {
		return nil, errors.New("El ID del libro es inválido")
	}

	_, err := s.bookStore.GetByID(libraryID, bookID)
	if err != nil {
		return nil, fmt.Errorf("El libro con ID %d no existe: %w", bookID, err)
	}

	categories, err := s.bookStore.GetBookCategories(libraryID, bookID)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las categorías del libro: %w", err)
	}

	return categories, nil
}

func (s *BookService) AddCategoryToBook(libraryID int64, bookCategory *models.BookCategory) error {
	if err := validations.ValidateBookCategory(bookCategory); err != nil {
		return fmt.Errorf("Validación fallida: %w", err)
	}

	_, err := s.bookStore.GetByID(libraryID, bookCategory.BookID)
	if err != nil {
		return fmt.Errorf("El libro con ID %d no existe: %w", bookCategory.BookID, err)
	}

	if err := s.bookStore.AddCategoryToBook(libraryID, bookCategory); err != nil {
		return fmt.Errorf("Error al agregar la categoría al libro: %w", err)
	}

	return nil
}

func (s *BookService) RemoveCategoryFromBook(libraryID, bookID, categoryID int64) error {
	if bookID <= 0 {
		return errors.New("El ID del libro es inválido")
	}

	if categoryID <= 0 {
		return errors.New("El ID de la categoría es inválido")
	}

	_, err := s.bookStore.GetByID(libraryID, bookID)
	if err != nil {
		return fmt.Errorf("El libro con ID %d no existe: %w", bookID, err)
	}

	if err := s.bookStore.RemoveCategoryFromBook(libraryID, bookID, categoryID); err != nil {
		return fmt.Errorf("Error al eliminar la categoría del libro: %w", err)
	}

	return nil
}
