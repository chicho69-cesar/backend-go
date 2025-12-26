package validations

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

var (
	isbnRegex = regexp.MustCompile(`^(?:ISBN(?:-1[03])?:? )?(?=[0-9X]{10}$|(?=(?:[0-9]+[- ]){3})[- 0-9X]{13}$|97[89][0-9]{10}$|(?=(?:[0-9]+[- ]){4})[- 0-9]{17}$)(?:97[89][- ]?)?[0-9]{1,5}[- ]?[0-9]+[- ]?[0-9]+[- ]?[0-9X]$`)

	validStatuses = map[string]bool{
		"Available":   true,
		"Borrowed":    true,
		"Reserved":    true,
		"Maintenance": true,
	}
)

func ValidateBook(book *models.Book) error {
	if book == nil {
		return errors.New("El libro no puede ser nulo")
	}

	if strings.TrimSpace(book.ISBN) == "" {
		return errors.New("El ISBN es requerido")
	}

	isbn := strings.ReplaceAll(book.ISBN, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")

	if len(isbn) != 10 && len(isbn) != 13 {
		return errors.New("El ISBN debe tener 10 o 13 dígitos")
	}

	if !isbnRegex.MatchString(book.ISBN) {
		return errors.New("El formato del ISBN es inválido")
	}

	if strings.TrimSpace(book.Title) == "" {
		return errors.New("El título es requerido")
	}

	if len(book.Title) < 1 {
		return errors.New("El título debe tener al menos 1 caracter")
	}

	if len(book.Title) > 255 {
		return errors.New("El título no puede exceder 255 caracteres")
	}

	if book.Subtitle.Valid && len(book.Subtitle.String) > 255 {
		return errors.New("El subtítulo no puede exceder 255 caracteres")
	}

	if book.Edition.Valid && len(book.Edition.String) > 100 {
		return errors.New("La edición no puede exceder 100 caracteres")
	}

	if book.Language.Valid {
		if len(book.Language.String) < 2 {
			return errors.New("El idioma debe tener al menos 2 caracteres")
		}

		if len(book.Language.String) > 50 {
			return errors.New("El idioma no puede exceder 50 caracteres")
		}
	}

	if book.PublicationYear.Valid {
		currentYear := int64(time.Now().Year())

		if book.PublicationYear.Int64 < 1000 {
			return errors.New("El año de publicación debe ser mayor a 1000")
		}

		if book.PublicationYear.Int64 > currentYear+1 {
			return errors.New("El año de publicación no puede ser mayor al año actual")
		}
	}

	if book.Pages.Valid {
		if book.Pages.Int64 < 1 {
			return errors.New("El número de páginas debe ser al menos 1")
		}

		if book.Pages.Int64 > 100000 {
			return errors.New("El número de páginas no puede exceder 100,000")
		}
	}

	if book.Synopsis.Valid && len(book.Synopsis.String) > 5000 {
		return errors.New("La sinopsis no puede exceder 5000 caracteres")
	}

	if book.PublisherID.Valid && book.PublisherID.Int64 <= 0 {
		return errors.New("El ID del editor debe ser un número positivo")
	}

	if book.ShelfID.Valid && book.ShelfID.Int64 <= 0 {
		return errors.New("El ID del estante debe ser valido")
	}

	if strings.TrimSpace(book.Status) == "" {
		return errors.New("El estado es requerido")
	}

	if !validStatuses[book.Status] {
		return errors.New("El estado debe ser: Available, Borrowed, Reserved o Maintenance")
	}

	return nil
}

func ValidateBookAuthor(bookAuthor *models.BookAuthor) error {
	if bookAuthor == nil {
		return errors.New("La relación libro-autor no puede ser nula")
	}

	if bookAuthor.BookID <= 0 {
		return errors.New("El ID del libro debe ser un número positivo")
	}

	if bookAuthor.AuthorID <= 0 {
		return errors.New("El ID del autor debe ser un número positivo")
	}

	if bookAuthor.Position < 1 {
		return errors.New("La posición del autor debe ser al menos 1")
	}

	if bookAuthor.Position > 50 {
		return errors.New("La posición del autor no puede exceder 50")
	}

	return nil
}

func ValidateBookCategory(bookCategory *models.BookCategory) error {
	if bookCategory == nil {
		return errors.New("La relación libro-categoría no puede ser nula")
	}

	if bookCategory.BookID <= 0 {
		return errors.New("El ID del libro debe ser un número positivo")
	}

	if bookCategory.CategoryID <= 0 {
		return errors.New("El ID de la categoría debe ser un número positivo")
	}

	return nil
}
