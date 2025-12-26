package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/validations"
)

type UserService struct {
	userStore        store.IUserStore
	loanStore        store.ILoanStore
	reservationStore store.IReservationStore
	fineStore        store.IFineStore
}

func NewUserService(userStore store.IUserStore, loanStore store.ILoanStore, reservationStore store.IReservationStore, fineStore store.IFineStore) *UserService {
	return &UserService{
		userStore:        userStore,
		loanStore:        loanStore,
		reservationStore: reservationStore,
		fineStore:        fineStore,
	}
}

func (s *UserService) GetAllUsers(libraryID int64) ([]*models.User, error) {
	users, err := s.userStore.GetAll(libraryID)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los usuarios: %w", err)
	}

	return users, nil
}

func (s *UserService) GetUserByID(libraryID, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("El ID del usuario es inválido")
	}

	user, err := s.userStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el usuario con ID %d: %w", id, err)
	}

	return user, nil
}

func (s *UserService) GetUsersFiltered(libraryID int64, filter store.UserFilter) ([]*models.User, error) {
	if filter.Code != "" {
		filter.Code = strings.TrimSpace(strings.ToUpper(filter.Code))
	}

	if filter.DNI != "" {
		filter.DNI = strings.TrimSpace(strings.ToUpper(filter.DNI))
	}

	if filter.UserType != "" {
		filter.UserType = strings.TrimSpace(filter.UserType)
	}

	if filter.Status != "" {
		filter.Status = strings.TrimSpace(filter.Status)
	}

	users, err := s.userStore.GetUsersFiltered(libraryID, filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los usuarios filtrados: %w", err)
	}

	return users, nil
}

func (s *UserService) CreateUser(libraryID int64, user *models.User) (*models.User, error) {
	if err := validations.ValidateUser(user); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	existingUser, _ := s.userStore.GetByCode(libraryID, user.Code)
	if existingUser != nil {
		return nil, fmt.Errorf("Ya existe un usuario con el código %s", user.Code)
	}

	userWithDNI, _ := s.userStore.GetByDNI(libraryID, user.DNI)
	if userWithDNI != nil {
		return nil, fmt.Errorf("Ya existe un usuario con el DNI %s", user.DNI)
	}

	user.LibraryID = libraryID
	user.Code = strings.TrimSpace(strings.ToUpper(user.Code))
	user.DNI = strings.TrimSpace(strings.ToUpper(user.DNI))
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	if user.Email.Valid {
		user.Email.String = strings.TrimSpace(strings.ToLower(user.Email.String))
	}

	if user.Phone.Valid {
		user.Phone.String = strings.TrimSpace(user.Phone.String)
	}

	if user.Address.Valid {
		user.Address.String = strings.TrimSpace(user.Address.String)
	}

	if user.RegistrationDate.IsZero() {
		user.RegistrationDate = time.Now()
	}

	if strings.TrimSpace(user.Status) == "" {
		user.Status = "Active"
	}

	createdUser, err := s.userStore.Create(libraryID, user)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el usuario: %w", err)
	}

	return createdUser, nil
}

func (s *UserService) UpdateUser(libraryID, id int64, user *models.User) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("El ID del usuario es inválido")
	}

	existingUser, err := s.userStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("El usuario con ID %d no existe: %w", id, err)
	}

	if existingUser == nil {
		return nil, fmt.Errorf("El usuario con ID %d no fue encontrado", id)
	}

	if err := validations.ValidateUser(user); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	userWithCode, _ := s.userStore.GetByCode(libraryID, user.Code)
	if userWithCode != nil && userWithCode.ID != id {
		return nil, fmt.Errorf("Ya existe otro usuario con el código %s", user.Code)
	}

	userWithDNI, _ := s.userStore.GetByDNI(libraryID, user.DNI)
	if userWithDNI != nil && userWithDNI.ID != id {
		return nil, fmt.Errorf("Ya existe otro usuario con el DNI %s", user.DNI)
	}

	user.Code = strings.TrimSpace(strings.ToUpper(user.Code))
	user.DNI = strings.TrimSpace(strings.ToUpper(user.DNI))
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	if user.Email.Valid {
		user.Email.String = strings.TrimSpace(strings.ToLower(user.Email.String))
	}

	if user.Phone.Valid {
		user.Phone.String = strings.TrimSpace(user.Phone.String)
	}

	if user.Address.Valid {
		user.Address.String = strings.TrimSpace(user.Address.String)
	}

	user.RegistrationDate = existingUser.RegistrationDate

	if user.Status == "Inactive" && existingUser.Status == "Active" {
		reservations, err := s.reservationStore.GetReservationsFiltered(libraryID, store.ReservationFilter{
			UserID: &id,
			Status: "Pending",
		})
		if err == nil && len(reservations) > 0 {
			for _, reservation := range reservations {
				reservation.Status = "Cancelled"
				s.reservationStore.Update(libraryID, reservation.ID, reservation)
			}
		}
	}

	updatedUser, err := s.userStore.Update(libraryID, id, user)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar el usuario con ID %d: %w", id, err)
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(libraryID, id int64) error {
	if id <= 0 {
		return errors.New("El ID del usuario es inválido")
	}

	existingUser, err := s.userStore.GetByID(libraryID, id)
	if err != nil {
		return fmt.Errorf("El usuario con ID %d no existe: %w", id, err)
	}

	if existingUser == nil {
		return fmt.Errorf("El usuario con ID %d no fue encontrado", id)
	}

	activeLoans, err := s.loanStore.GetLoansFiltered(libraryID, store.LoanFilter{
		UserID: &id,
		Status: "Active",
	})
	if err == nil && len(activeLoans) > 0 {
		return fmt.Errorf("No se puede eliminar el usuario porque tiene %d préstamo(s) activo(s)", len(activeLoans))
	}

	overdueLoans, err := s.loanStore.GetLoansFiltered(libraryID, store.LoanFilter{
		UserID:  &id,
		Overdue: true,
	})
	if err == nil && len(overdueLoans) > 0 {
		return fmt.Errorf("No se puede eliminar el usuario porque tiene %d préstamo(s) vencido(s)", len(overdueLoans))
	}

	pendingFines, err := s.fineStore.GetFinesFiltered(libraryID, store.FineFilter{
		UserID:  &id,
		Pending: true,
	})
	if err == nil && len(pendingFines) > 0 {
		return fmt.Errorf("No se puede eliminar el usuario porque tiene %d multa(s) pendiente(s)", len(pendingFines))
	}

	userReservations, err := s.reservationStore.GetReservationsFiltered(libraryID, store.ReservationFilter{
		UserID: &id,
	})
	if err == nil && len(userReservations) > 0 {
		for _, reservation := range userReservations {
			s.reservationStore.Delete(libraryID, reservation.ID)
		}
	}

	if err := s.userStore.Delete(libraryID, id); err != nil {
		return fmt.Errorf("Error al eliminar el usuario con ID %d: %w", id, err)
	}

	return nil
}
