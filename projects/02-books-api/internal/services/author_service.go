package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/validations"
)

type AuthorService struct {
	authorStore store.AuthStore
}

func NewAuthorService(authorStore store.AuthStore) *AuthorService {
	return &AuthorService{
		authorStore: authorStore,
	}
}

func (s *AuthorService) GetAllAuthors() ([]*models.Author, error) {
	return s.authorStore.GetAll()
}

func (s *AuthorService) GetAuthorByID(id int64) (*models.Author, error) {
	if id <= 0 {
		return nil, errors.New("El ID del autor es invalido")
	}

	author, err := s.authorStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el autor con ID %d: %w", id, err)
	}

	return author, nil
}

func (s *AuthorService) CreateAuthor(author *models.Author) (*models.Author, error) {
	if err := validations.ValidateAuthor(author); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	author.FirstName = strings.TrimSpace(author.FirstName)
	author.LastName = strings.TrimSpace(author.LastName)

	if author.Biography.Valid {
		author.Biography.String = strings.TrimSpace(author.Biography.String)
	}

	if author.Nationality.Valid {
		author.Nationality.String = strings.TrimSpace(author.Nationality.String)
	}

	createdAuthor, err := s.authorStore.Create(author)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el autor: %w", err)
	}

	return createdAuthor, nil
}

func (s *AuthorService) UpdateAuthor(id int64, author *models.Author) (*models.Author, error) {
	if id <= 0 {
		return nil, errors.New("El ID del autor es invalido")
	}

	existingAuthor, err := s.authorStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("El autor con ID %d no existe: %w", id, err)
	}

	if existingAuthor == nil {
		return nil, fmt.Errorf("El autor con ID %d no fue encontrado", id)
	}

	if err := validations.ValidateAuthor(author); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	author.FirstName = strings.TrimSpace(author.FirstName)
	author.LastName = strings.TrimSpace(author.LastName)

	if author.Biography.Valid {
		author.Biography.String = strings.TrimSpace(author.Biography.String)
	}

	if author.Nationality.Valid {
		author.Nationality.String = strings.TrimSpace(author.Nationality.String)
	}

	updatedAuthor, err := s.authorStore.Update(id, author)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar el autor con ID %d: %w", id, err)
	}

	return updatedAuthor, nil
}

func (s *AuthorService) DeleteAuthor(id int64) error {
	if id <= 0 {
		return errors.New("El ID del autor es invalido")
	}

	existingAuthor, err := s.authorStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("El autor con ID %d no existe: %w", id, err)
	}

	if existingAuthor == nil {
		return fmt.Errorf("El autor con ID %d no fue encontrado", id)
	}

	if err := s.authorStore.Delete(id); err != nil {
		return fmt.Errorf("Error al eliminar el autor con ID %d: %w", id, err)
	}

	return nil
}
