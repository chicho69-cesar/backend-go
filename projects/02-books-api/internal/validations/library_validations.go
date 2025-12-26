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
