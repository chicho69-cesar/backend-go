package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/middleware"
	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
)

type LoanHandler struct {
	loanService *services.LoanService
}

type ReservationHandler struct {
	reservationService *services.ReservationService
}

type FineHandler struct {
	fineService *services.FineService
}

func NewLoanHandler(loanService *services.LoanService) *LoanHandler {
	return &LoanHandler{loanService: loanService}
}

func NewReservationHandler(reservationService *services.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func NewFineHandler(fineService *services.FineService) *FineHandler {
	return &FineHandler{fineService: fineService}
}

// GET /loans - Obtener todos los préstamos o con filtros
// POST /loans - Crear un nuevo préstamo
func (h *LoanHandler) HandleLoans(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.LoanFilter{}

			code := r.URL.Query().Get("code")
			if code != "" {
				filter.Code = code
				hasFilters = true
			}

			userIDStr := r.URL.Query().Get("user_id")
			if userIDStr != "" {
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil || userID <= 0 {
					http.Error(w, "El parámetro user_id es inválido", http.StatusBadRequest)
					return
				}
				
				filter.UserID = &userID
				hasFilters = true
			}

			copyIDStr := r.URL.Query().Get("copy_id")
			if copyIDStr != "" {
				copyID, err := strconv.ParseInt(copyIDStr, 10, 64)
				if err != nil || copyID <= 0 {
					http.Error(w, "El parámetro copy_id es inválido", http.StatusBadRequest)
					return
				}

				filter.CopyID = &copyID
				hasFilters = true
			}

			status := r.URL.Query().Get("status")
			if status != "" {
				filter.Status = status
				hasFilters = true
			}

			overdue := r.URL.Query().Get("overdue")
			if overdue == "true" {
				filter.Overdue = true
				hasFilters = true
			}

			var loans []*models.Loan
			var err error

			if hasFilters {
				loans, err = h.loanService.GetLoansFiltered(libraryID, filter)
			} else {
				loans, err = h.loanService.GetAll(libraryID)
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("Error al obtener préstamos: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(loans)

		case http.MethodPost:
			var loan models.Loan
			if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
				http.Error(w, fmt.Sprintf("Error al decodificar body: %v", err), http.StatusBadRequest)
				return
			}

			createdLoan, err := h.loanService.CreateLoan(libraryID, &loan)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al crear préstamo: %v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(createdLoan)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /loans/{id} - Obtener préstamo por ID
// PUT /loans/{id} - Actualizar préstamo por ID
// DELETE /loans/{id} - Eliminar préstamo por ID
func (h *LoanHandler) HandleLoanByID(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			loan, err := h.loanService.GetByID(libraryID, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al obtener préstamo: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(loan)

		case http.MethodPut:
			var loan models.Loan
			if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
				http.Error(w, fmt.Sprintf("Error al decodificar body: %v", err), http.StatusBadRequest)
				return
			}

			updatedLoan, err := h.loanService.Update(libraryID, id, &loan)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al actualizar préstamo: %v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedLoan)

		case http.MethodDelete:
			if err := h.loanService.Delete(libraryID, id); err != nil {
				http.Error(w, fmt.Sprintf("Error al eliminar préstamo: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// POST /loans/{id}/renew - Renovar préstamo
func (h *LoanHandler) HandleLoanRenew(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var body map[string]any
	var librarianID *int64

	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		if libID, ok := body["librarian_id"].(float64); ok {
			libIDInt := int64(libID)
			librarianID = &libIDInt
		}
	}

	renewedLoan, err := h.loanService.RenewLoan(libraryID, id, librarianID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al renovar préstamo: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(renewedLoan)
}

// POST /loans/{id}/return - Devolver préstamo
func (h *LoanHandler) HandleLoanReturn(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var body map[string]interface{}
	var notes *string

	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		if notesStr, ok := body["notes"].(string); ok {
			notes = &notesStr
		}
	}

	returnedLoan, err := h.loanService.ReturnLoan(libraryID, id, notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al devolver préstamo: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnedLoan)
}

// GET /reservations - Obtener todas las reservaciones o con filtros
// POST /reservations - Crear una nueva reservación
func (h *ReservationHandler) HandleReservations(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.ReservationFilter{}

			userIDStr := r.URL.Query().Get("user_id")
			if userIDStr != "" {
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil || userID <= 0 {
					http.Error(w, "El parámetro user_id es inválido", http.StatusBadRequest)
					return
				}

				filter.UserID = &userID
				hasFilters = true
			}

			bookIDStr := r.URL.Query().Get("book_id")
			if bookIDStr != "" {
				bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
				if err != nil || bookID <= 0 {
					http.Error(w, "El parámetro book_id es inválido", http.StatusBadRequest)
					return
				}

				filter.BookID = &bookID
				hasFilters = true
			}

			status := r.URL.Query().Get("status")
			if status != "" {
				filter.Status = status
				hasFilters = true
			}

			expired := r.URL.Query().Get("expired")
			if expired == "true" {
				filter.Expired = true
				hasFilters = true
			}

			var reservations []*models.Reservation
			var err error

			if hasFilters {
				reservations, err = h.reservationService.GetReservationsFiltered(libraryID, filter)
			} else {
				reservations, err = h.reservationService.GetAll(libraryID)
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("Error al obtener reservaciones: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(reservations)

		case http.MethodPost:
			var reservation models.Reservation
			if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
				http.Error(w, fmt.Sprintf("Error al decodificar body: %v", err), http.StatusBadRequest)
				return
			}

			if reservation.ReservationDate.IsZero() {
				reservation.ReservationDate = time.Now()
			}

			if reservation.ExpirationDate.IsZero() {
				reservation.ExpirationDate = time.Now().Add(7 * 24 * time.Hour)
			}

			if reservation.Status == "" {
				reservation.Status = "Pending"
			}

			if reservation.Priority == 0 {
				reservation.Priority = 5
			}

			createdReservation, err := h.reservationService.CreateReservation(libraryID, &reservation)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al crear reservación: %v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(createdReservation)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /reservations/{id} - Obtener reservación por ID
// PUT /reservations/{id} - Actualizar reservación por ID
// DELETE /reservations/{id} - Eliminar reservación por ID
func (h *ReservationHandler) HandleReservationByID(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			reservation, err := h.reservationService.GetByID(libraryID, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al obtener reservación: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(reservation)

		case http.MethodPut:
			var reservation models.Reservation
			if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
				http.Error(w, fmt.Sprintf("Error al decodificar body: %v", err), http.StatusBadRequest)
				return
			}

			updatedReservation, err := h.reservationService.Update(libraryID, id, &reservation)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al actualizar reservación: %v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedReservation)

		case http.MethodDelete:
			if err := h.reservationService.Delete(libraryID, id); err != nil {
				http.Error(w, fmt.Sprintf("Error al eliminar reservación: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// POST /reservations/{id}/cancel - Cancelar reservación
func (h *ReservationHandler) HandleReservationCancel(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	cancelledReservation, err := h.reservationService.CancelReservation(libraryID, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al cancelar reservación: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cancelledReservation)
}

// POST /reservations/{id}/process - Procesar reservación
func (h *ReservationHandler) HandleReservationProcess(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	processedReservation, err := h.reservationService.ProcessReservation(libraryID, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al procesar reservación: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processedReservation)
}

// GET /fines - Obtener todas las multas o con filtros
// POST /fines - Crear una nueva multa
func (h *FineHandler) HandleFines(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.FineFilter{}

			userIDStr := r.URL.Query().Get("user_id")
			if userIDStr != "" {
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil || userID <= 0 {
					http.Error(w, "El parámetro user_id es inválido", http.StatusBadRequest)
					return
				}

				filter.UserID = &userID
				hasFilters = true
			}

			loanIDStr := r.URL.Query().Get("loan_id")
			if loanIDStr != "" {
				loanID, err := strconv.ParseInt(loanIDStr, 10, 64)
				if err != nil || loanID <= 0 {
					http.Error(w, "El parámetro loan_id es inválido", http.StatusBadRequest)
					return
				}

				filter.LoanID = &loanID
				hasFilters = true
			}

			status := r.URL.Query().Get("status")
			if status != "" {
				filter.Status = status
				hasFilters = true
			}

			pending := r.URL.Query().Get("pending")
			if pending == "true" {
				filter.Pending = true
				hasFilters = true
			}

			var fines []*models.Fine
			var err error

			if hasFilters {
				fines, err = h.fineService.GetFinesFiltered(libraryID, filter)
			} else {
				fines, err = h.fineService.GetAll(libraryID)
			}

			if err != nil {
				http.Error(w, fmt.Sprintf("Error al obtener multas: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fines)

		case http.MethodPost:
			var fine models.Fine
			if err := json.NewDecoder(r.Body).Decode(&fine); err != nil {
				http.Error(w, fmt.Sprintf("Error al decodificar body: %v", err), http.StatusBadRequest)
				return
			}

			if fine.GeneratedDate.IsZero() {
				fine.GeneratedDate = time.Now()
			}

			if fine.Status == "" {
				fine.Status = "Pending"
			}

			createdFine, err := h.fineService.CreateFine(libraryID, &fine)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al crear multa: %v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(createdFine)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /fines/{id} - Obtener multa por ID
// PUT /fines/{id} - Actualizar multa por ID
// DELETE /fines/{id} - Eliminar multa por ID
func (h *FineHandler) HandleFineByID(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			fine, err := h.fineService.GetByID(libraryID, id)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al obtener multa: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(fine)

		case http.MethodPut:
			var fine models.Fine
			if err := json.NewDecoder(r.Body).Decode(&fine); err != nil {
				http.Error(w, fmt.Sprintf("Error al decodificar body: %v", err), http.StatusBadRequest)
				return
			}

			updatedFine, err := h.fineService.Update(libraryID, id, &fine)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error al actualizar multa: %v", err), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedFine)

		case http.MethodDelete:
			if err := h.fineService.Delete(libraryID, id); err != nil {
				http.Error(w, fmt.Sprintf("Error al eliminar multa: %v", err), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// POST /fines/{id}/pay - Pagar multa
func (h *FineHandler) HandleFinePay(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var body map[string]any
	var notes *string

	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		if notesStr, ok := body["notes"].(string); ok {
			notes = &notesStr
		}
	}

	paidFine, err := h.fineService.PayFine(libraryID, id, notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al pagar multa: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paidFine)
}

// POST /fines/{id}/waive - Condonar multa
func (h *FineHandler) HandleFineWaive(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var body map[string]any
	var notes *string

	if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
		if notesStr, ok := body["notes"].(string); ok {
			notes = &notesStr
		}
	}

	waivedFine, err := h.fineService.WaiveFine(libraryID, id, notes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al condonar multa: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(waivedFine)
}
