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

type LibraryService struct {
	libraryStore store.ILibraryStore
}

type LibraryZoneService struct {
	zoneStore store.ILibraryZoneStore
}

type ShelfService struct {
	shelfStore store.IShelfStore
	zoneStore  store.ILibraryZoneStore
}

type CopyService struct {
	copyStore store.ICopyStore
	bookStore store.IBookStore
	loanStore store.ILoanStore
}

func NewLibraryService(libraryStore store.ILibraryStore) *LibraryService {
	return &LibraryService{libraryStore: libraryStore}
}

func NewLibraryZoneService(zoneStore store.ILibraryZoneStore) *LibraryZoneService {
	return &LibraryZoneService{zoneStore: zoneStore}
}

func NewShelfService(shelfStore store.IShelfStore, zoneStore store.ILibraryZoneStore) *ShelfService {
	return &ShelfService{
		shelfStore: shelfStore,
		zoneStore:  zoneStore,
	}
}

func NewCopyService(copyStore store.ICopyStore, bookStore store.IBookStore, loanStore store.ILoanStore) *CopyService {
	return &CopyService{
		copyStore: copyStore,
		bookStore: bookStore,
		loanStore: loanStore,
	}
}

func (s *LibraryService) GetAllLibraries() ([]*models.Library, error) {
	libraries, err := s.libraryStore.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la información: %w", err)
	}

	return libraries, nil
}

func (s *LibraryService) GetLibraryByID(id int64) (*models.Library, error) {
	if id <= 0 {
		return nil, errors.New("El ID de la biblioteca es inválido")
	}

	library, err := s.libraryStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la biblioteca con ID %d: %w", id, err)
	}

	return library, nil
}

func (s *LibraryService) EnterLibraryCredentials(username, password string) (*models.Library, error) {
	library, err := s.libraryStore.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la biblioteca con el nombre de usuario %s: %w", username, err)
	}

	if library == nil {
		return nil, errors.New("Nombre de usuario o contraseña incorrectos")
	}

	checkedWithPassword, err := s.libraryStore.CheckPassword(username, password)
	if err != nil {
		return nil, errors.New("Nombre de usuario o contraseña incorrectos")
	}

	if !checkedWithPassword {
		return nil, errors.New("Nombre de usuario o contraseña incorrectos")
	}

	return library, nil
}

func (s *LibraryService) CreateLibrary(library *models.Library) (*models.Library, error) {
	if err := validations.ValidateLibrary(library); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	existingLibrary, _ := s.libraryStore.GetByUsername(library.Username)
	if existingLibrary != nil {
		return nil, fmt.Errorf("Ya existe una biblioteca con el nombre de usuario %s", library.Username)
	}

	createdLibrary, err := s.libraryStore.Create(library)
	if err != nil {
		return nil, fmt.Errorf("Error al crear la biblioteca: %w", err)
	}

	return createdLibrary, nil
}

func (s *LibraryService) UpdateLibrary(id int64, library *models.Library) (*models.Library, error) {
	if id <= 0 {
		return nil, errors.New("El ID de la biblioteca es inválido")
	}

	existingLibrary, err := s.libraryStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("La biblioteca con ID %d no existe: %w", id, err)
	}

	if existingLibrary == nil {
		return nil, fmt.Errorf("La biblioteca con ID %d no fue encontrada", id)
	}

	if err := validations.ValidateLibrary(library); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	updatedLibrary, err := s.libraryStore.Update(id, library)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar la biblioteca con ID %d: %w", id, err)
	}

	return updatedLibrary, nil
}

func (s *LibraryService) DeleteLibrary(id int64) error {
	if id <= 0 {
		return errors.New("El ID de la biblioteca es inválido")
	}

	existingLibrary, err := s.libraryStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("La biblioteca con ID %d no existe: %w", id, err)
	}

	if existingLibrary == nil {
		return fmt.Errorf("La biblioteca con ID %d no fue encontrada", id)
	}

	if err := s.libraryStore.Delete(id); err != nil {
		return fmt.Errorf("Error al eliminar la biblioteca con ID %d: %w", id, err)
	}

	return nil
}

func (s *LibraryZoneService) GetAllZones() ([]*models.LibraryZone, error) {
	zones, err := s.zoneStore.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las zonas: %w", err)
	}

	return zones, nil
}

func (s *LibraryZoneService) GetZoneByID(id int64) (*models.LibraryZone, error) {
	if id <= 0 {
		return nil, errors.New("El ID de la zona es inválido")
	}

	zone, err := s.zoneStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la zona con ID %d: %w", id, err)
	}

	return zone, nil
}

func (s *LibraryZoneService) GetZonesFiltered(filter store.LibraryZoneFilter) ([]*models.LibraryZone, error) {
	if filter.Code != "" {
		filter.Code = strings.TrimSpace(strings.ToUpper(filter.Code))
	}

	if filter.Floor != nil {
		if *filter.Floor < 0 {
			return nil, errors.New("El piso debe ser 0 o mayor")
		}
	}

	zones, err := s.zoneStore.GetZonesFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las zonas filtradas: %w", err)
	}

	return zones, nil
}

func (s *LibraryZoneService) CreateZone(zone *models.LibraryZone) (*models.LibraryZone, error) {
	if err := validations.ValidateLibraryZone(zone); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	existingZone, _ := s.zoneStore.GetByCode(zone.Code)
	if existingZone != nil {
		return nil, fmt.Errorf("Ya existe una zona con el código %s", zone.Code)
	}

	zone.Code = strings.TrimSpace(strings.ToUpper(zone.Code))
	zone.Name = strings.TrimSpace(zone.Name)

	if zone.Description.Valid {
		zone.Description.String = strings.TrimSpace(zone.Description.String)
	}

	createdZone, err := s.zoneStore.Create(zone)
	if err != nil {
		return nil, fmt.Errorf("Error al crear la zona: %w", err)
	}

	return createdZone, nil
}

func (s *LibraryZoneService) UpdateZone(id int64, zone *models.LibraryZone) (*models.LibraryZone, error) {
	if id <= 0 {
		return nil, errors.New("El ID de la zona es inválido")
	}

	existingZone, err := s.zoneStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("La zona con ID %d no existe: %w", id, err)
	}

	if existingZone == nil {
		return nil, fmt.Errorf("La zona con ID %d no fue encontrada", id)
	}

	if err := validations.ValidateLibraryZone(zone); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	zoneWithCode, _ := s.zoneStore.GetByCode(zone.Code)
	if zoneWithCode != nil && zoneWithCode.ID != id {
		return nil, fmt.Errorf("Ya existe otra zona con el código %s", zone.Code)
	}

	zone.Code = strings.TrimSpace(strings.ToUpper(zone.Code))
	zone.Name = strings.TrimSpace(zone.Name)

	if zone.Description.Valid {
		zone.Description.String = strings.TrimSpace(zone.Description.String)
	}

	updatedZone, err := s.zoneStore.Update(id, zone)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar la zona con ID %d: %w", id, err)
	}

	return updatedZone, nil
}

func (s *LibraryZoneService) DeleteZone(id int64) error {
	if id <= 0 {
		return errors.New("El ID de la zona es inválido")
	}

	existingZone, err := s.zoneStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("La zona con ID %d no existe: %w", id, err)
	}

	if existingZone == nil {
		return fmt.Errorf("La zona con ID %d no fue encontrada", id)
	}

	if err := s.zoneStore.Delete(id); err != nil {
		return fmt.Errorf("Error al eliminar la zona con ID %d: %w", id, err)
	}

	return nil
}

func (s *ShelfService) GetAllShelves() ([]*models.Shelf, error) {
	shelves, err := s.shelfStore.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los estantes: %w", err)
	}

	return shelves, nil
}

func (s *ShelfService) GetShelfByID(id int64) (*models.Shelf, error) {
	if id <= 0 {
		return nil, errors.New("El ID del estante es inválido")
	}

	shelf, err := s.shelfStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener el estante con ID %d: %w", id, err)
	}

	return shelf, nil
}

func (s *ShelfService) GetShelvesFiltered(filter store.ShelfFilter) ([]*models.Shelf, error) {
	if filter.Code != "" {
		filter.Code = strings.TrimSpace(strings.ToUpper(filter.Code))
	}

	if filter.ZoneID != nil {
		if *filter.ZoneID <= 0 {
			return nil, errors.New("El ID de la zona es inválido")
		}

		_, err := s.zoneStore.GetByID(*filter.ZoneID)
		if err != nil {
			return nil, fmt.Errorf("La zona con ID %d no existe: %w", *filter.ZoneID, err)
		}
	}

	shelves, err := s.shelfStore.GetShelvesFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener los estantes filtrados: %w", err)
	}

	return shelves, nil
}

func (s *ShelfService) CreateShelf(shelf *models.Shelf) (*models.Shelf, error) {
	if err := validations.ValidateShelf(shelf); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	_, err := s.zoneStore.GetByID(shelf.ZoneID)
	if err != nil {
		return nil, fmt.Errorf("La zona con ID %d no existe: %w", shelf.ZoneID, err)
	}

	existingShelf, _ := s.shelfStore.GetByCode(shelf.Code)
	if existingShelf != nil {
		return nil, fmt.Errorf("Ya existe un estante con el código %s", shelf.Code)
	}

	shelf.Code = strings.TrimSpace(strings.ToUpper(shelf.Code))

	if shelf.Description.Valid {
		shelf.Description.String = strings.TrimSpace(shelf.Description.String)
	}

	createdShelf, err := s.shelfStore.Create(shelf)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el estante: %w", err)
	}

	return createdShelf, nil
}

func (s *ShelfService) UpdateShelf(id int64, shelf *models.Shelf) (*models.Shelf, error) {
	if id <= 0 {
		return nil, errors.New("El ID del estante es inválido")
	}

	existingShelf, err := s.shelfStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("El estante con ID %d no existe: %w", id, err)
	}

	if existingShelf == nil {
		return nil, fmt.Errorf("El estante con ID %d no fue encontrado", id)
	}

	if err := validations.ValidateShelf(shelf); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	_, err = s.zoneStore.GetByID(shelf.ZoneID)
	if err != nil {
		return nil, fmt.Errorf("La zona con ID %d no existe: %w", shelf.ZoneID, err)
	}

	shelfWithCode, _ := s.shelfStore.GetByCode(shelf.Code)
	if shelfWithCode != nil && shelfWithCode.ID != id {
		return nil, fmt.Errorf("Ya existe otro estante con el código %s", shelf.Code)
	}

	shelf.Code = strings.TrimSpace(strings.ToUpper(shelf.Code))

	if shelf.Description.Valid {
		shelf.Description.String = strings.TrimSpace(shelf.Description.String)
	}

	updatedShelf, err := s.shelfStore.Update(id, shelf)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar el estante con ID %d: %w", id, err)
	}

	return updatedShelf, nil
}

func (s *ShelfService) DeleteShelf(id int64) error {
	if id <= 0 {
		return errors.New("El ID del estante es inválido")
	}

	existingShelf, err := s.shelfStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("El estante con ID %d no existe: %w", id, err)
	}

	if existingShelf == nil {
		return fmt.Errorf("El estante con ID %d no fue encontrado", id)
	}

	if err := s.shelfStore.Delete(id); err != nil {
		return fmt.Errorf("Error al eliminar el estante con ID %d: %w", id, err)
	}

	return nil
}

func (s *CopyService) GetAllCopies() ([]*models.Copy, error) {
	copies, err := s.copyStore.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las copias: %w", err)
	}

	return copies, nil
}

func (s *CopyService) GetCopyByID(id int64) (*models.Copy, error) {
	if id <= 0 {
		return nil, errors.New("El ID de la copia es inválido")
	}

	copy, err := s.copyStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la copia con ID %d: %w", id, err)
	}

	return copy, nil
}

func (s *CopyService) GetCopiesFiltered(filter store.CopyFilter) ([]*models.Copy, error) {
	if filter.Code != "" {
		filter.Code = strings.TrimSpace(strings.ToUpper(filter.Code))
	}

	if filter.BookID != nil {
		if *filter.BookID <= 0 {
			return nil, errors.New("El ID del libro es inválido")
		}

		_, err := s.bookStore.GetByID(*filter.BookID)
		if err != nil {
			return nil, fmt.Errorf("El libro con ID %d no existe: %w", *filter.BookID, err)
		}
	}

	if filter.Status != "" {
		filter.Status = strings.TrimSpace(filter.Status)
	}

	if filter.Condition != "" {
		filter.Condition = strings.TrimSpace(filter.Condition)
	}

	copies, err := s.copyStore.GetCopiesFiltered(filter)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener las copias filtradas: %w", err)
	}

	return copies, nil
}

func (s *CopyService) CreateCopy(copy *models.Copy) (*models.Copy, error) {
	if err := validations.ValidateCopy(copy); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	_, err := s.bookStore.GetByID(copy.BookID)
	if err != nil {
		return nil, fmt.Errorf("El libro con ID %d no existe: %w", copy.BookID, err)
	}

	existingCopy, _ := s.copyStore.GetByCode(copy.Code)
	if existingCopy != nil {
		return nil, fmt.Errorf("Ya existe una copia con el código %s", copy.Code)
	}

	copy.Code = strings.TrimSpace(strings.ToUpper(copy.Code))

	if copy.Notes.Valid {
		copy.Notes.String = strings.TrimSpace(copy.Notes.String)
	}

	if !copy.AcquisitionDate.Valid {
		copy.AcquisitionDate.Time = time.Now()
		copy.AcquisitionDate.Valid = true
	}

	createdCopy, err := s.copyStore.Create(copy)
	if err != nil {
		return nil, fmt.Errorf("Error al crear la copia: %w", err)
	}

	return createdCopy, nil
}

func (s *CopyService) UpdateCopy(id int64, copy *models.Copy) (*models.Copy, error) {
	if id <= 0 {
		return nil, errors.New("El ID de la copia es inválido")
	}

	existingCopy, err := s.copyStore.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("La copia con ID %d no existe: %w", id, err)
	}

	if existingCopy == nil {
		return nil, fmt.Errorf("La copia con ID %d no fue encontrada", id)
	}

	if err := validations.ValidateCopy(copy); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	_, err = s.bookStore.GetByID(copy.BookID)
	if err != nil {
		return nil, fmt.Errorf("El libro con ID %d no existe: %w", copy.BookID, err)
	}

	copyWithCode, _ := s.copyStore.GetByCode(copy.Code)
	if copyWithCode != nil && copyWithCode.ID != id {
		return nil, fmt.Errorf("Ya existe otra copia con el código %s", copy.Code)
	}

	if copy.Status != "Borrowed" && existingCopy.Status == "Borrowed" {
		activeLoans, err := s.loanStore.GetLoansFiltered(store.LoanFilter{
			CopyID: &id,
			Status: "Active",
		})
		if err == nil && len(activeLoans) > 0 {
			return nil, fmt.Errorf("No se puede cambiar el estado porque la copia tiene %d préstamo(s) activo(s)", len(activeLoans))
		}
	}

	copy.Code = strings.TrimSpace(strings.ToUpper(copy.Code))

	if copy.Notes.Valid {
		copy.Notes.String = strings.TrimSpace(copy.Notes.String)
	}

	updatedCopy, err := s.copyStore.Update(id, copy)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar la copia con ID %d: %w", id, err)
	}

	return updatedCopy, nil
}

func (s *CopyService) DeleteCopy(id int64) error {
	if id <= 0 {
		return errors.New("El ID de la copia es inválido")
	}

	existingCopy, err := s.copyStore.GetByID(id)
	if err != nil {
		return fmt.Errorf("La copia con ID %d no existe: %w", id, err)
	}

	if existingCopy == nil {
		return fmt.Errorf("La copia con ID %d no fue encontrada", id)
	}

	if existingCopy.Status == "Borrowed" {
		return fmt.Errorf("No se puede eliminar la copia porque está actualmente prestada")
	}

	activeLoans, err := s.loanStore.GetLoansFiltered(store.LoanFilter{
		CopyID: &id,
		Status: "Active",
	})
	if err == nil && len(activeLoans) > 0 {
		return fmt.Errorf("No se puede eliminar la copia porque tiene %d préstamo(s) activo(s)", len(activeLoans))
	}

	overdueLoans, err := s.loanStore.GetLoansFiltered(store.LoanFilter{
		CopyID:  &id,
		Overdue: true,
	})
	if err == nil && len(overdueLoans) > 0 {
		return fmt.Errorf("No se puede eliminar la copia porque tiene %d préstamo(s) vencido(s) sin devolver", len(overdueLoans))
	}

	if err := s.copyStore.Delete(id); err != nil {
		return fmt.Errorf("Error al eliminar la copia con ID %d: %w", id, err)
	}

	return nil
}
