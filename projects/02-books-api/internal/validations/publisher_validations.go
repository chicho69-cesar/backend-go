package validations

import (
	"errors"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

func ValidatePublisher(publisher *models.Publisher) error {
	if publisher == nil {
		return errors.New("La editorial no puede estar vacía")
	}

	if strings.TrimSpace(publisher.Name) == "" {
		return errors.New("El nombre de la editorial es requerido")
	}

	if len(publisher.Name) < 2 {
		return errors.New("El nombre de la editorial debe tener al menos 2 caracteres")
	}

	if len(publisher.Name) > 100 {
		return errors.New("El nombre de la editorial no puede exceder 100 caracteres")
	}

	if publisher.Country.Valid {
		if len(publisher.Country.String) < 2 {
			return errors.New("El país de la editorial debe tener al menos 2 caracteres")
		}

		if len(publisher.Country.String) > 100 {
			return errors.New("El país de la editorial no puede exceder 100 caracteres")
		}
	}

	return nil
}
