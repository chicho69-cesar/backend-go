package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/validations"
)

type LoanService struct {
	loanStore store.ILoanStore
	userStore store.IUserStore
	copyStore store.ICopyStore
	fineStore store.IFineStore
}

type ReservationService struct {
	reservationStore store.IReservationStore
	userStore        store.IUserStore
	bookStore        store.IBookStore
	copyStore        store.ICopyStore
	fineStore        store.IFineStore
}

type FineService struct {
	fineStore store.IFineStore
	userStore store.IUserStore
	loanStore store.ILoanStore
}

func NewLoanService(loanStore store.ILoanStore, userStore store.IUserStore, copyStore store.ICopyStore, fineStore store.IFineStore) *LoanService {
	return &LoanService{
		loanStore: loanStore,
		userStore: userStore,
		copyStore: copyStore,
		fineStore: fineStore,
	}
}

func NewReservationService(reservationStore store.IReservationStore, userStore store.IUserStore, bookStore store.IBookStore, copyStore store.ICopyStore, fineStore store.IFineStore) *ReservationService {
	return &ReservationService{
		reservationStore: reservationStore,
		userStore:        userStore,
		bookStore:        bookStore,
		copyStore:        copyStore,
		fineStore:        fineStore,
	}
}

func NewFineService(fineStore store.IFineStore, userStore store.IUserStore, loanStore store.ILoanStore) *FineService {
	return &FineService{
		fineStore: fineStore,
		userStore: userStore,
		loanStore: loanStore,
	}
}

func (s *LoanService) GetAll() ([]*models.Loan, error) {
	return s.loanStore.GetAll()
}

func (s *LoanService) GetByID(id int64) (*models.Loan, error) {
	return s.loanStore.GetByID(id)
}

func (s *LoanService) GetByCode(code string) (*models.Loan, error) {
	return s.loanStore.GetByCode(code)
}

func (s *LoanService) GetLoansFiltered(filter store.LoanFilter) ([]*models.Loan, error) {
	if filter.Code != "" {
		filter.Code = strings.TrimSpace(strings.ToUpper(filter.Code))
	}

	if filter.Status != "" {
		filter.Status = strings.TrimSpace(filter.Status)
	}

	loans, err := s.loanStore.GetLoansFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los préstamos filtrados: %w", err)
	}

	return loans, nil
}

func (s *LoanService) CreateLoan(loan *models.Loan) (*models.Loan, error) {
	if err := validations.ValidateLoan(loan); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetByID(loan.UserID)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar usuario: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("El usuario con ID %d no existe", loan.UserID)
	}

	if user.Status != "Active" {
		return nil, fmt.Errorf("El usuario no está activo")
	}

	copy, err := s.copyStore.GetByID(loan.CopyID)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar copia: %v", err)
	}

	if copy == nil {
		return nil, fmt.Errorf("La copia con ID %d no existe", loan.CopyID)
	}

	if copy.Status != "Available" {
		return nil, fmt.Errorf("La copia no está disponible para préstamo")
	}

	if copy.Condition == "Damaged" || copy.Condition == "Lost" {
		return nil, fmt.Errorf("La copia no está en condiciones para préstamo")
	}

	overdueLoans, err := s.loanStore.GetLoansFiltered(store.LoanFilter{Overdue: true})
	if err != nil {
		return nil, fmt.Errorf("Error al verificar préstamos vencidos: %v", err)
	}

	for _, overdueLoan := range overdueLoans {
		if overdueLoan.UserID == loan.UserID {
			return nil, fmt.Errorf("El usuario tiene préstamos vencidos, no puede realizar nuevos préstamos")
		}
	}

	userID := loan.UserID
	pendingFines, err := s.fineStore.GetFinesFiltered(store.FineFilter{UserID: &userID, Pending: true})
	if err != nil {
		return nil, fmt.Errorf("Error al verificar multas pendientes: %v", err)
	}

	if len(pendingFines) > 0 {
		return nil, fmt.Errorf("El usuario tiene multas pendientes, debe pagarlas antes de realizar un préstamo")
	}

	existingLoan, err := s.loanStore.GetByCode(loan.LoanCode)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar código de préstamo: %v", err)
	}

	if existingLoan != nil {
		return nil, fmt.Errorf("El código de préstamo %s ya existe", loan.LoanCode)
	}

	createdLoan, err := s.loanStore.Create(loan)
	if err != nil {
		return nil, fmt.Errorf("Error al crear préstamo: %v", err)
	}

	copy.Status = "Loaned"
	_, err = s.copyStore.Update(copy.ID, copy)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar estado de copia: %v", err)
	}

	return createdLoan, nil
}

func (s *LoanService) RenewLoan(id int64, librarianID *int64) (*models.Loan, error) {
	loan, err := s.loanStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener préstamo: %v", err)
	}

	if loan == nil {
		return nil, fmt.Errorf("Préstamo con ID %d no encontrado", id)
	}

	if err := validations.ValidateLoanRenewal(loan); err != nil {
		return nil, err
	}

	userID := loan.UserID
	pendingFines, err := s.fineStore.GetFinesFiltered(store.FineFilter{UserID: &userID, Pending: true})
	if err != nil {
		return nil, fmt.Errorf("Error al verificar multas pendientes: %v", err)
	}

	if len(pendingFines) > 0 {
		return nil, fmt.Errorf("El usuario tiene multas pendientes, debe pagarlas antes de renovar el préstamo")
	}

	loan.Renewals++
	loan.DueDate = loan.DueDate.Add(time.Duration(loan.LoanDays) * 24 * time.Hour)

	if librarianID != nil {
		loan.LibrarianID.Valid = true
		loan.LibrarianID.Int64 = *librarianID
	}

	updatedLoan, err := s.loanStore.Update(id, loan)
	if err != nil {
		return nil, fmt.Errorf("Error al renovar préstamo: %v", err)
	}

	return updatedLoan, nil
}

func (s *LoanService) ReturnLoan(id int64, notes *string) (*models.Loan, error) {
	loan, err := s.loanStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener préstamo: %v", err)
	}

	if loan == nil {
		return nil, fmt.Errorf("Préstamo con ID %d no encontrado", id)
	}

	if loan.Status != "Active" {
		return nil, fmt.Errorf("El préstamo no está activo")
	}

	now := time.Now()
	loan.Status = "Returned"
	loan.ReturnDate.Valid = true
	loan.ReturnDate.Time = now

	if notes != nil {
		loan.Notes.Valid = true
		loan.Notes.String = *notes
	}

	updatedLoan, err := s.loanStore.Update(id, loan)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar préstamo: %v", err)
	}

	copy, err := s.copyStore.GetByID(loan.CopyID)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener copia: %v", err)
	}

	copy.Status = "Available"
	_, err = s.copyStore.Update(copy.ID, copy)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar estado de copia: %v", err)
	}

	if now.After(loan.DueDate) {
		daysLate := int(now.Sub(loan.DueDate).Hours() / 24)
		fineAmount := float64(daysLate) * 5.0

		if fineAmount > 0 {
			fine := &models.Fine{
				UserID:        loan.UserID,
				Reason:        fmt.Sprintf("Devolución tardía (%d días)", daysLate),
				Amount:        fineAmount,
				GeneratedDate: now,
				Status:        "Pending",
			}

			fine.LoanID.Valid = true
			fine.LoanID.Int64 = loan.ID

			_, err = s.fineStore.Create(fine)
			if err != nil {
				fmt.Printf("Advertencia: Error al crear multa automática: %v\n", err)
			}
		}
	}

	return updatedLoan, nil
}

func (s *LoanService) Update(id int64, loan *models.Loan) (*models.Loan, error) {
	if err := validations.ValidateLoan(loan); err != nil {
		return nil, err
	}

	existingLoan, err := s.loanStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener préstamo: %v", err)
	}

	if existingLoan == nil {
		return nil, fmt.Errorf("Préstamo con ID %d no encontrado", id)
	}

	return s.loanStore.Update(id, loan)
}

func (s *LoanService) Delete(id int64) error {
	loan, err := s.loanStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("Error al obtener préstamo: %v", err)
	}

	if loan == nil {
		return fmt.Errorf("Préstamo con ID %d no encontrado", id)
	}

	return s.loanStore.Delete(id)
}

func (s *ReservationService) GetAll() ([]*models.Reservation, error) {
	return s.reservationStore.GetAll()
}

func (s *ReservationService) GetByID(id int64) (*models.Reservation, error) {
	return s.reservationStore.GetByID(id)
}

func (s *ReservationService) GetReservationsFiltered(filter store.ReservationFilter) ([]*models.Reservation, error) {
	if filter.Status != "" {
		filter.Status = strings.TrimSpace(filter.Status)
	}

	reservations, err := s.reservationStore.GetReservationsFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las reservaciones filtradas: %w", err)
	}

	return reservations, nil
}

func (s *ReservationService) CreateReservation(reservation *models.Reservation) (*models.Reservation, error) {
	if err := validations.ValidateReservation(reservation); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetByID(reservation.UserID)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar usuario: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("El usuario con ID %d no existe", reservation.UserID)
	}

	if user.Status != "Active" {
		return nil, fmt.Errorf("El usuario no está activo")
	}

	pendingFines, err := s.fineStore.GetFinesFiltered(store.FineFilter{
		UserID:  &reservation.UserID,
		Pending: true,
	})
	
	if err == nil && len(pendingFines) > 0 {
		return nil, fmt.Errorf("El usuario tiene %d multa(s) pendiente(s) y no puede hacer reservaciones", len(pendingFines))
	}

	book, err := s.bookStore.GetByID(reservation.BookID)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar libro: %v", err)
	}

	if book == nil {
		return nil, fmt.Errorf("El libro con ID %d no existe", reservation.BookID)
	}

	existingReservation, err := s.reservationStore.GetActiveByUserAndBook(reservation.UserID, reservation.BookID)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar reservaciones existentes: %v", err)
	}

	if existingReservation != nil {
		return nil, fmt.Errorf("El usuario ya tiene una reservación activa para este libro")
	}

	bookID := reservation.BookID
	copies, err := s.copyStore.GetCopiesFiltered(store.CopyFilter{BookID: &bookID})
	if err != nil {
		return nil, fmt.Errorf("Error al verificar copias del libro: %v", err)
	}

	availableCopies := 0
	for _, copy := range copies {
		if copy.Status == "Available" {
			availableCopies++
		}
	}

	if availableCopies > 0 {
		return nil, fmt.Errorf("Hay copias disponibles del libro, no es necesario realizar una reservación")
	}

	createdReservation, err := s.reservationStore.Create(reservation)
	if err != nil {
		return nil, fmt.Errorf("Error al crear reservación: %v", err)
	}

	return createdReservation, nil
}

func (s *ReservationService) CancelReservation(id int64) (*models.Reservation, error) {
	reservation, err := s.reservationStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener reservación: %v", err)
	}

	if reservation == nil {
		return nil, fmt.Errorf("Reservación con ID %d no encontrada", id)
	}

	if reservation.Status != "Pending" && reservation.Status != "Active" {
		return nil, fmt.Errorf("La reservación no se puede cancelar (estado actual: %s)", reservation.Status)
	}

	reservation.Status = "Cancelled"

	updatedReservation, err := s.reservationStore.Update(id, reservation)
	if err != nil {
		return nil, fmt.Errorf("Error al cancelar reservación: %v", err)
	}

	return updatedReservation, nil
}

func (s *ReservationService) ProcessReservation(id int64) (*models.Reservation, error) {
	reservation, err := s.reservationStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener reservación: %v", err)
	}

	if reservation == nil {
		return nil, fmt.Errorf("Reservación con ID %d no encontrada", id)
	}

	if reservation.Status != "Active" {
		return nil, fmt.Errorf("La reservación no está activa")
	}

	reservation.Status = "Completed"

	updatedReservation, err := s.reservationStore.Update(id, reservation)
	if err != nil {
		return nil, fmt.Errorf("Error al procesar reservación: %v", err)
	}

	return updatedReservation, nil
}

func (s *ReservationService) Update(id int64, reservation *models.Reservation) (*models.Reservation, error) {
	if err := validations.ValidateReservation(reservation); err != nil {
		return nil, err
	}

	existingReservation, err := s.reservationStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener reservación: %v", err)
	}

	if existingReservation == nil {
		return nil, fmt.Errorf("Reservación con ID %d no encontrada", id)
	}

	return s.reservationStore.Update(id, reservation)
}

func (s *ReservationService) Delete(id int64) error {
	reservation, err := s.reservationStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("Error al obtener reservación: %v", err)
	}

	if reservation == nil {
		return fmt.Errorf("Reservación con ID %d no encontrada", id)
	}

	return s.reservationStore.Delete(id)
}

func (s *FineService) GetAll() ([]*models.Fine, error) {
	return s.fineStore.GetAll()
}

func (s *FineService) GetByID(id int64) (*models.Fine, error) {
	return s.fineStore.GetByID(id)
}

func (s *FineService) GetFinesFiltered(filter store.FineFilter) ([]*models.Fine, error) {
	if filter.Status != "" {
		filter.Status = strings.TrimSpace(filter.Status)
	}

	fines, err := s.fineStore.GetFinesFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las multas filtradas: %w", err)
	}

	return fines, nil
}

func (s *FineService) CreateFine(fine *models.Fine) (*models.Fine, error) {
	if err := validations.ValidateFine(fine); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetByID(fine.UserID)
	if err != nil {
		return nil, fmt.Errorf("Error al verificar usuario: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("El usuario con ID %d no existe", fine.UserID)
	}

	if fine.LoanID.Valid {
		loan, err := s.loanStore.GetByID(fine.LoanID.Int64)
		if err != nil {
			return nil, fmt.Errorf("Error al verificar préstamo: %v", err)
		}

		if loan == nil {
			return nil, fmt.Errorf("El préstamo con ID %d no existe", fine.LoanID.Int64)
		}
	}

	createdFine, err := s.fineStore.Create(fine)
	if err != nil {
		return nil, fmt.Errorf("Error al crear multa: %v", err)
	}

	return createdFine, nil
}

func (s *FineService) PayFine(id int64, notes *string) (*models.Fine, error) {
	fine, err := s.fineStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener multa: %v", err)
	}

	if fine == nil {
		return nil, fmt.Errorf("Multa con ID %d no encontrada", id)
	}

	if fine.Status != "Pending" {
		return nil, fmt.Errorf("La multa no está pendiente de pago")
	}

	now := time.Now()
	fine.Status = "Paid"
	fine.PaymentDate.Valid = true
	fine.PaymentDate.Time = now

	if notes != nil {
		fine.Notes.Valid = true
		fine.Notes.String = *notes
	}

	updatedFine, err := s.fineStore.Update(id, fine)
	if err != nil {
		return nil, fmt.Errorf("Error al marcar multa como pagada: %v", err)
	}

	return updatedFine, nil
}

func (s *FineService) WaiveFine(id int64, notes *string) (*models.Fine, error) {
	fine, err := s.fineStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener multa: %v", err)
	}

	if fine == nil {
		return nil, fmt.Errorf("Multa con ID %d no encontrada", id)
	}

	if fine.Status != "Pending" {
		return nil, fmt.Errorf("La multa no está pendiente")
	}

	fine.Status = "Waived"

	if notes != nil {
		fine.Notes.Valid = true
		fine.Notes.String = *notes
	}

	updatedFine, err := s.fineStore.Update(id, fine)
	if err != nil {
		return nil, fmt.Errorf("Error al condonar multa: %v", err)
	}

	return updatedFine, nil
}

func (s *FineService) Update(id int64, fine *models.Fine) (*models.Fine, error) {
	if err := validations.ValidateFine(fine); err != nil {
		return nil, err
	}

	existingFine, err := s.fineStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener multa: %v", err)
	}

	if existingFine == nil {
		return nil, fmt.Errorf("Multa con ID %d no encontrada", id)
	}

	return s.fineStore.Update(id, fine)
}

func (s *FineService) Delete(id int64) error {
	fine, err := s.fineStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("Error al obtener multa: %v", err)
	}

	if fine == nil {
		return fmt.Errorf("Multa con ID %d no encontrada", id)
	}

	return s.fineStore.Delete(id)
}
