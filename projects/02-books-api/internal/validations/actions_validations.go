package validations

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
)

var (
	loanCodeRegex = regexp.MustCompile(`^LOAN-\d{4}-\d{4,6}$`)

	validLoanStatuses = map[string]bool{
		"Active":   true,
		"Returned": true,
		"Overdue":  true,
		"Lost":     true,
	}

	validReservationStatuses = map[string]bool{
		"Pending":   true,
		"Active":    true,
		"Cancelled": true,
		"Expired":   true,
	}

	validFineReasons = map[string]bool{
		"Overdue": true,
		"Damage":  true,
		"Loss":    true,
	}

	validFineStatuses = map[string]bool{
		"Pending": true,
		"Paid":    true,
		"Waived":  true,
	}
)

func ValidateLoan(loan *models.Loan) error {
	if loan == nil {
		return errors.New("El préstamo no puede ser nulo")
	}

	if strings.TrimSpace(loan.LoanCode) == "" {
		return errors.New("El código de préstamo es requerido")
	}

	if !loanCodeRegex.MatchString(loan.LoanCode) {
		return errors.New("El código debe tener el formato LOAN-YYYY-NNNN (ej: LOAN-2024-0001)")
	}

	if loan.UserID <= 0 {
		return errors.New("El ID del usuario debe ser un número positivo")
	}

	if loan.CopyID <= 0 {
		return errors.New("El ID de la copia debe ser un número positivo")
	}

	if loan.LoanDate.IsZero() {
		return errors.New("La fecha de préstamo es requerida")
	}

	if loan.DueDate.IsZero() {
		return errors.New("La fecha de vencimiento es requerida")
	}

	if loan.DueDate.Before(loan.LoanDate) {
		return errors.New("La fecha de vencimiento debe ser posterior a la fecha de préstamo")
	}

	if loan.ReturnDate.Valid {
		if loan.ReturnDate.Time.Before(loan.LoanDate) {
			return errors.New("La fecha de devolución no puede ser anterior a la fecha de préstamo")
		}
	}

	if strings.TrimSpace(loan.Status) == "" {
		return errors.New("El estado es requerido")
	}

	if !validLoanStatuses[loan.Status] {
		return errors.New("El estado debe ser: Active, Returned, Overdue o Lost")
	}

	if loan.LoanDays < 1 {
		return errors.New("Los días de préstamo deben ser al menos 1")
	}

	if loan.LoanDays > 90 {
		return errors.New("Los días de préstamo no pueden exceder 90")
	}

	if loan.Renewals < 0 {
		return errors.New("Las renovaciones no pueden ser negativas")
	}

	if loan.Renewals > 5 {
		return errors.New("No se pueden hacer más de 5 renovaciones")
	}

	if loan.Notes.Valid && len(loan.Notes.String) > 1000 {
		return errors.New("Las notas no pueden exceder 1000 caracteres")
	}

	if loan.LibrarianID.Valid && loan.LibrarianID.Int64 <= 0 {
		return errors.New("El ID del bibliotecario debe ser un número positivo")
	}

	return nil
}

func ValidateReservation(reservation *models.Reservation) error {
	if reservation == nil {
		return errors.New("La reservación no puede ser nula")
	}

	if reservation.UserID <= 0 {
		return errors.New("El ID del usuario debe ser un número positivo")
	}

	if reservation.BookID <= 0 {
		return errors.New("El ID del libro debe ser un número positivo")
	}

	if reservation.ReservationDate.IsZero() {
		return errors.New("La fecha de reservación es requerida")
	}

	if reservation.ExpirationDate.IsZero() {
		return errors.New("La fecha de expiración es requerida")
	}

	if reservation.ExpirationDate.Before(reservation.ReservationDate) {
		return errors.New("La fecha de expiración debe ser posterior a la fecha de reservación")
	}

	if strings.TrimSpace(reservation.Status) == "" {
		return errors.New("El estado es requerido")
	}

	if !validReservationStatuses[reservation.Status] {
		return errors.New("El estado debe ser: Pending, Active, Cancelled o Expired")
	}

	if reservation.Priority < 1 {
		return errors.New("La prioridad debe ser al menos 1")
	}

	if reservation.Priority > 10 {
		return errors.New("La prioridad no puede exceder 10")
	}

	return nil
}

func ValidateFine(fine *models.Fine) error {
	if fine == nil {
		return errors.New("La multa no puede ser nula")
	}

	if fine.UserID <= 0 {
		return errors.New("El ID del usuario debe ser un número positivo")
	}

	if fine.LoanID.Valid && fine.LoanID.Int64 <= 0 {
		return errors.New("El ID del préstamo debe ser un número positivo")
	}

	if strings.TrimSpace(fine.Reason) == "" {
		return errors.New("La razón es requerida")
	}

	if !validFineReasons[fine.Reason] {
		return errors.New("La razón debe ser: Overdue, Damage o Loss")
	}

	if fine.Amount < 0 {
		return errors.New("El monto no puede ser negativo")
	}

	if fine.Amount > 10000 {
		return errors.New("El monto no puede exceder 10,000")
	}

	if fine.GeneratedDate.IsZero() {
		return errors.New("La fecha de generación es requerida")
	}

	if fine.PaymentDate.Valid {
		if fine.PaymentDate.Time.Before(fine.GeneratedDate) {
			return errors.New("La fecha de pago no puede ser anterior a la fecha de generación")
		}
	}

	if strings.TrimSpace(fine.Status) == "" {
		return errors.New("El estado es requerido")
	}

	if !validFineStatuses[fine.Status] {
		return errors.New("El estado debe ser: Pending, Paid o Waived")
	}

	if fine.Notes.Valid && len(fine.Notes.String) > 1000 {
		return errors.New("Las notas no pueden exceder 1000 caracteres")
	}

	return nil
}

func ValidateLoanRenewal(loan *models.Loan) error {
	if loan.Status != "Active" {
		return errors.New("Solo se pueden renovar préstamos activos")
	}

	if loan.Renewals >= 5 {
		return errors.New("Se ha alcanzado el máximo de renovaciones permitidas")
	}

	if loan.Status == "Overdue" {
		daysPastDue := int(time.Since(loan.DueDate).Hours() / 24)
		if daysPastDue > 3 {
			return errors.New("No se puede renovar un préstamo vencido por más de 3 días")
		}
	}

	return nil
}
