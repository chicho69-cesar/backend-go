package validations

import (
	"errors"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

func ValidateCategory(category *models.Category) error {
	if category == nil {
		return errors.New("La categoría no puede estar vacía")
	}

	if strings.TrimSpace(category.Name) == "" {
		return errors.New("El nombre de la categoría es requerido")
	}

	if len(category.Name) < 2 {
		return errors.New("El nombre de la categoría debe tener al menos 2 caracteres")
	}

	if len(category.Name) > 100 {
		return errors.New("El nombre de la categoría no puede exceder 100 caracteres")
	}

	if category.Description.Valid {
		if len(category.Description.String) < 5 {
			return errors.New("La descripción de la categoría debe tener al menos 5 caracteres")
		}

		if len(category.Description.String) > 500 {
			return errors.New("La descripción de la categoría no puede exceder 500 caracteres")
		}
	}

	return nil
}