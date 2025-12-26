package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/validations"
)

type CategoryService struct {
	categoryStore store.ICategoryStore
}

func NewCategoryService(categoryStore store.ICategoryStore) *CategoryService {
	return &CategoryService{
		categoryStore: categoryStore,
	}
}

func (s *CategoryService) GetAllCategories(libraryID int64) ([]*models.Category, error) {
	return s.categoryStore.GetAll(libraryID)
}

func (s *CategoryService) GetCategoryByID(libraryID, id int64) (*models.Category, error) {
	if id <= 0 {
		return nil, errors.New("El id de la categoría es invalido")
	}

	category, err := s.categoryStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la categoría con ID %d: %w", id, err)
	}

	return category, nil
}

func (s *CategoryService) CreateCategory(libraryID int64, category *models.Category) (*models.Category, error) {
	if err := validations.ValidateCategory(category); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	category.LibraryID = libraryID
	category.Name = strings.TrimSpace(category.Name)

	if category.Description.Valid {
		category.Description.String = strings.TrimSpace(category.Description.String)
	}

	createdCategory, err := s.categoryStore.Create(libraryID, category)
	if err != nil {
		return nil, fmt.Errorf("Error al crear la categoría: %w", err)
	}

	return createdCategory, nil
}

func (s *CategoryService) UpdateCategory(libraryID, id int64, category *models.Category) (*models.Category, error) {
	if id <= 0 {
		return nil, errors.New("El id de la categoría es invalido")
	}

	existingCategory, err := s.categoryStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la categoría con ID %d: %w", id, err)
	}

	if existingCategory == nil {
		return nil, fmt.Errorf("La categoría con ID %d no existe", id)
	}

	if err := validations.ValidateCategory(category); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	category.Name = strings.TrimSpace(category.Name)

	if category.Description.Valid {
		category.Description.String = strings.TrimSpace(category.Description.String)
	}

	updatedCategory, err := s.categoryStore.Update(libraryID, id, category)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar la categoría con ID %d: %w", id, err)
	}

	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategory(libraryID, id int64) error {
	if id <= 0 {
		return errors.New("El id de la categoría es invalido")
	}

	existingCategory, err := s.categoryStore.GetByID(libraryID, id)
	if err != nil {
		return fmt.Errorf("Error al obtener la categoría con ID %d: %w", id, err)
	}

	if existingCategory == nil {
		return fmt.Errorf("La categoría con ID %d no existe", id)
	}

	if err := s.categoryStore.Delete(libraryID, id); err != nil {
		return fmt.Errorf("Error al eliminar la categoría con ID %d: %w", id, err)
	}

	return nil
}
