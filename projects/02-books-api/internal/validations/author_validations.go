package validations

import (
	"errors"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

func ValidateAuthor(author *models.Author) error {
	if author == nil {
		return errors.New("El autor no puede ser nulo")
	}

	if strings.TrimSpace(author.FirstName) == "" {
		return errors.New("El nombre es requerido")
	}

	if len(author.FirstName) < 2 {
		return errors.New("El nombre debe tener al menos 2 caracteres")
	}

	if len(author.FirstName) > 100 {
		return errors.New("El nombre no puede exceder 100 caracteres")
	}

	if strings.TrimSpace(author.LastName) == "" {
		return errors.New("El apellido es requerido")
	}

	if len(author.LastName) < 2 {
		return errors.New("El apellido debe tener al menos 2 caracteres")
	}

	if len(author.LastName) > 100 {
		return errors.New("El apellido no puede exceder 100 caracteres")
	}

	if author.Nationality.Valid {
		if len(author.Nationality.String) < 2 {
			return errors.New("La nacionalidad debe tener al menos 2 caracteres")
		}

		if len(author.Nationality.String) > 200 {
			return errors.New("La nacionalidad no puede exceder 200 caracteres")
		}
	}

	return nil
}
