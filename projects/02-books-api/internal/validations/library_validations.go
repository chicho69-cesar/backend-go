package validations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

var (
	zoneCodeRegex  = regexp.MustCompile(`^[A-Z][0-9]{1,2}$`)
	shelfCodeRegex = regexp.MustCompile(`^[A-Z][0-9]{1,2}-[0-9]{2}$`)
	copyCodeRegex  = regexp.MustCompile(`^[A-Z0-9]{6,20}$`)
	websiteRegex   = regexp.MustCompile(`^(https?://)?(www\.)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[\w.-]*)*/?$`)
	passwordRegex  = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{6,100}$`)

	validCopyStatuses = map[string]bool{
		"Available": true,
		"Borrowed":  true,
		"Reserved":  true,
		"Damaged":   true,
		"Lost":      true,
	}

	validCopyConditions = map[string]bool{
		"New":  true,
		"Good": true,
		"Fair": true,
		"Poor": true,
	}
)

func ValidateLibrary(library *models.Library) error {
	if library == nil {
		return errors.New("La biblioteca no puede estar vacía")
	}

	if strings.TrimSpace(library.Name) == "" {
		return errors.New("El nombre de la biblioteca es requerido")
	}

	if len(library.Name) < 2 {
		return errors.New("El nombre de la biblioteca debe tener al menos 2 caracteres")
	}

	if len(library.Name) > 255 {
		return errors.New("El nombre de la biblioteca no puede exceder 255 caracteres")
	}

	if strings.TrimSpace(library.Address) == "" {
		return errors.New("La dirección es requerida")
	}

	if strings.TrimSpace(library.City) == "" {
		return errors.New("La ciudad es requerida")
	}

	if strings.TrimSpace(library.Country) == "" {
		return errors.New("El país es requerido")
	}

	if strings.TrimSpace(library.Email) == "" {
		return errors.New("El correo electrónico es requerido")
	}

	if !emailRegex.MatchString(library.Email) {
		return errors.New("El formato del email es inválido")
	}

	if len(library.Email) > 255 {
		return errors.New("El email no puede exceder 255 caracteres")
	}

	if strings.TrimSpace(library.Phone) == "" {
		return errors.New("El teléfono es requerido")
	}

	cleanPhone := strings.ReplaceAll(library.Phone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")

	if !phoneRegex.MatchString(cleanPhone) {
		return errors.New("El teléfono debe tener entre 10-15 dígitos")
	}

	if strings.TrimSpace(library.State) == "" {
		return errors.New("El estado/provincia es requerido")
	}

	if strings.TrimSpace(library.Username) == "" {
		return errors.New("El nombre de usuario es requerido")
	}

	if strings.TrimSpace(library.Website) != "" {
		if !websiteRegex.MatchString(library.Website) {
			return errors.New("El formato del sitio web es inválido")
		}
	}

	if strings.TrimSpace(library.Password) == "" {
		return errors.New("La contraseña es requerida")
	}

	if strings.TrimSpace(library.Password) != "" {
		if len(library.Password) < 6 {
			return errors.New("La contraseña debe tener al menos 6 caracteres")
		}

		if len(library.Password) > 100 {
			return errors.New("La contraseña no puede exceder 100 caracteres")
		}

		if !passwordRegex.MatchString(library.Password) {
			return errors.New("La contraseña debe contener al menos una letra mayúscula, una letra minúscula, un número y un carácter especial")
		}
	}

	if strings.TrimSpace(library.ZipCode) == "" {
		return errors.New("El código postal es requerido")
	}

	return nil
}

func ValidateLibraryZone(zone *models.LibraryZone) error {
	if zone == nil {
		return errors.New("La zona no puede estar vacía")
	}

	if strings.TrimSpace(zone.Code) == "" {
		return errors.New("El código de zona es requerido")
	}

	if !zoneCodeRegex.MatchString(zone.Code) {
		return errors.New("El código de zona debe tener el formato: letra mayúscula seguida de 1-2 dígitos (ej: A1, B12)")
	}

	if strings.TrimSpace(zone.Name) == "" {
		return errors.New("El nombre de la zona es requerido")
	}

	if len(zone.Name) < 2 {
		return errors.New("El nombre de la zona debe tener al menos 2 caracteres")
	}

	if len(zone.Name) > 100 {
		return errors.New("El nombre de la zona no puede exceder 100 caracteres")
	}

	if zone.Description.Valid && len(zone.Description.String) > 500 {
		return errors.New("La descripción no puede exceder 500 caracteres")
	}

	if zone.Floor < 0 {
		return errors.New("El piso debe ser 0 o mayor")
	}

	if zone.Floor > 50 {
		return errors.New("El piso no puede exceder 50")
	}

	return nil
}

func ValidateShelf(shelf *models.Shelf) error {
	if shelf == nil {
		return errors.New("El estante no puede estar vacío")
	}

	if strings.TrimSpace(shelf.Code) == "" {
		return errors.New("El código del estante es requerido")
	}

	if !shelfCodeRegex.MatchString(shelf.Code) {
		return errors.New("El código del estante debe tener el formato: zona-número (ej: A1-01, B2-15)")
	}

	if shelf.ZoneID <= 0 {
		return errors.New("El ID de la zona debe ser un número positivo")
	}

	if shelf.Description.Valid && len(shelf.Description.String) > 500 {
		return errors.New("La descripción no puede exceder 500 caracteres")
	}

	return nil
}

func ValidateCopy(copy *models.Copy) error {
	if copy == nil {
		return errors.New("La copia no puede estar vacía")
	}

	if strings.TrimSpace(copy.Code) == "" {
		return errors.New("El código/barcode de la copia es requerido")
	}

	if !copyCodeRegex.MatchString(copy.Code) {
		return errors.New("El código debe ser alfanumérico de 6-20 caracteres")
	}

	if copy.BookID <= 0 {
		return errors.New("El ID del libro debe ser un número positivo")
	}

	if strings.TrimSpace(copy.Status) == "" {
		return errors.New("El estado es requerido")
	}

	if !validCopyStatuses[copy.Status] {
		return errors.New("El estado debe ser: Available, Borrowed, Reserved, Damaged o Lost")
	}

	if strings.TrimSpace(copy.Condition) == "" {
		return errors.New("La condición es requerida")
	}

	if !validCopyConditions[copy.Condition] {
		return errors.New("La condición debe ser: New, Good, Fair o Poor")
	}

	if copy.PurchasePrice.Valid && copy.PurchasePrice.Float64 < 0 {
		return errors.New("El precio de compra no puede ser negativo")
	}

	if copy.Notes.Valid && len(copy.Notes.String) > 1000 {
		return errors.New("Las notas no pueden exceder 1000 caracteres")
	}

	return nil
}
