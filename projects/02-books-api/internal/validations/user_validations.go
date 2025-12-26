package validations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

var (
	userCodeRegex = regexp.MustCompile(`^[A-Z0-9]{4,20}$`)
	dniRegex      = regexp.MustCompile(`^[A-Z0-9]{6,15}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^[+]?[0-9]{10,15}$`)

	validUserTypes = map[string]bool{
		"Student":  true,
		"Teacher":  true,
		"Staff":    true,
		"External": true,
	}

	validUserStatuses = map[string]bool{
		"Active":    true,
		"Suspended": true,
		"Inactive":  true,
	}
)

func ValidateUser(user *models.User) error {
	if user == nil {
		return errors.New("El usuario no puede ser nulo")
	}

	if strings.TrimSpace(user.Code) == "" {
		return errors.New("El código de usuario es requerido")
	}

	if !userCodeRegex.MatchString(user.Code) {
		return errors.New("El código debe ser alfanumérico de 4-20 caracteres")
	}

	if strings.TrimSpace(user.DNI) == "" {
		return errors.New("El DNI/documento es requerido")
	}

	if !dniRegex.MatchString(user.DNI) {
		return errors.New("El DNI debe ser alfanumérico de 6-15 caracteres")
	}

	if strings.TrimSpace(user.FirstName) == "" {
		return errors.New("El nombre es requerido")
	}

	if len(user.FirstName) < 2 {
		return errors.New("El nombre debe tener al menos 2 caracteres")
	}

	if len(user.FirstName) > 100 {
		return errors.New("El nombre no puede exceder 100 caracteres")
	}

	if strings.TrimSpace(user.LastName) == "" {
		return errors.New("El apellido es requerido")
	}

	if len(user.LastName) < 2 {
		return errors.New("El apellido debe tener al menos 2 caracteres")
	}

	if len(user.LastName) > 100 {
		return errors.New("El apellido no puede exceder 100 caracteres")
	}

	if user.Email.Valid {
		if !emailRegex.MatchString(user.Email.String) {
			return errors.New("El formato del email es inválido")
		}

		if len(user.Email.String) > 255 {
			return errors.New("El email no puede exceder 255 caracteres")
		}
	}

	if user.Phone.Valid {
		cleanPhone := strings.ReplaceAll(user.Phone.String, "-", "")
		cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")

		if !phoneRegex.MatchString(cleanPhone) {
			return errors.New("El teléfono debe tener entre 10-15 dígitos")
		}
	}

	if user.Address.Valid && len(user.Address.String) > 500 {
		return errors.New("La dirección no puede exceder 500 caracteres")
	}

	if strings.TrimSpace(user.UserType) == "" {
		return errors.New("El tipo de usuario es requerido")
	}

	if !validUserTypes[user.UserType] {
		return errors.New("El tipo de usuario debe ser: Student, Teacher, Staff o External")
	}

	if strings.TrimSpace(user.Status) == "" {
		return errors.New("El estado es requerido")
	}

	if !validUserStatuses[user.Status] {
		return errors.New("El estado debe ser: Active, Suspended o Inactive")
	}

	return nil
}
