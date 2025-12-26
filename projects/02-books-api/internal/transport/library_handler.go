package transport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
)

type LibraryHandler struct {
	libraryService *services.LibraryService
}

type LibraryZoneHandler struct {
	zoneService *services.LibraryZoneService
}

type ShelfHandler struct {
	shelfService *services.ShelfService
}

type CopyHandler struct {
	copyService *services.CopyService
}

func NewLibraryHandler(libraryService *services.LibraryService) *LibraryHandler {
	return &LibraryHandler{libraryService: libraryService}
}

func NewLibraryZoneHandler(zoneService *services.LibraryZoneService) *LibraryZoneHandler {
	return &LibraryZoneHandler{zoneService: zoneService}
}

func NewShelfHandler(shelfService *services.ShelfService) *ShelfHandler {
	return &ShelfHandler{shelfService: shelfService}
}

func NewCopyHandler(copyService *services.CopyService) *CopyHandler {
	return &CopyHandler{copyService: copyService}
}

// GET /libraries - Obtener todas las bibliotecas
// POST /libraries - Crear una nueva biblioteca
func (h *LibraryHandler) HandleLibraries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			libraries, err := h.libraryService.GetAllLibraries()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(libraries)
		
		case http.MethodPost:
			var library models.Library
			err := json.NewDecoder(r.Body).Decode(&library)
			if err != nil {
				http.Error(w, "Datos de biblioteca inválidos", http.StatusBadRequest)
				return
			}

			createdLibrary, err := h.libraryService.CreateLibrary(&library)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdLibrary)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /libraries/{id} - Obtener biblioteca por ID
// POST /libraries/{id} - Ingresar a la biblioteca con credenciales
// PUT /libraries/{id} - Actualizar biblioteca por ID
// DELETE /libraries/{id} - Eliminar biblioteca por ID
func (h *LibraryHandler) HandleLibraryByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/libraries/"):]
	if idParam == "" {
		http.Error(w, "El parámetro ID es requerido", http.StatusBadRequest)
		return
	}

	readId, err := strconv.Atoi(idParam)
	if err != nil || readId <= 0 {
		http.Error(w, "El ID es inválido", http.StatusBadRequest)
		return
	}

	id := int64(readId)

	switch r.Method {
		case http.MethodGet:
			library, err := h.libraryService.GetLibraryByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(library)
		
		case http.MethodPost:
			var credentials struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			err := json.NewDecoder(r.Body).Decode(&credentials)
			if err != nil {
				http.Error(w, "Credenciales inválidas", http.StatusBadRequest)
				return
			}

			_, err = h.libraryService.GetLibraryByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			authenticatedLibrary, err := h.libraryService.EnterLibraryCredentials(credentials.Username, credentials.Password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(authenticatedLibrary)

		case http.MethodPut:
			var library models.Library
			err := json.NewDecoder(r.Body).Decode(&library)
			if err != nil {
				http.Error(w, "Datos de biblioteca inválidos", http.StatusBadRequest)
				return
			}

			updatedLibrary, err := h.libraryService.UpdateLibrary(id, &library)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			json.NewEncoder(w).Encode(updatedLibrary)

		case http.MethodDelete:
			err := h.libraryService.DeleteLibrary(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /zones - Obtener todas las zonas o con filtros
// POST /zones - Crear una nueva zona
func (h *LibraryZoneHandler) HandleZones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.LibraryZoneFilter{}

			code := r.URL.Query().Get("code")
			if code != "" {
				filter.Code = code
				hasFilters = true
			}

			floorStr := r.URL.Query().Get("floor")
			if floorStr != "" {
				floor, err := strconv.Atoi(floorStr)
				if err != nil || floor < 0 {
					http.Error(w, "El parámetro floor es inválido", http.StatusBadRequest)
					return
				}

				filter.Floor = &floor
				hasFilters = true
			}

			var zones []*models.LibraryZone
			var err error

			if hasFilters {
				zones, err = h.zoneService.GetZonesFiltered(filter)
			} else {
				zones, err = h.zoneService.GetAllZones()
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(zones)

		case http.MethodPost:
			var zone models.LibraryZone
			err := json.NewDecoder(r.Body).Decode(&zone)
			if err != nil {
				http.Error(w, "Datos de zona inválidos", http.StatusBadRequest)
				return
			}

			createdZone, err := h.zoneService.CreateZone(&zone)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdZone)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /zones/{id} - Obtener zona por ID
// PUT /zones/{id} - Actualizar zona por ID
// DELETE /zones/{id} - Eliminar zona por ID
func (h *LibraryZoneHandler) HandleZoneByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/zones/")
	if idParam == "" {
		http.Error(w, "El parámetro ID es requerido", http.StatusBadRequest)
		return
	}

	readId, err := strconv.Atoi(idParam)
	if err != nil || readId <= 0 {
		http.Error(w, "El ID es inválido", http.StatusBadRequest)
		return
	}

	id := int64(readId)

	switch r.Method {
		case http.MethodGet:
			zone, err := h.zoneService.GetZoneByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(zone)

		case http.MethodPut:
			var zone models.LibraryZone
			err := json.NewDecoder(r.Body).Decode(&zone)
			if err != nil {
				http.Error(w, "Datos de zona inválidos", http.StatusBadRequest)
				return
			}

			updatedZone, err := h.zoneService.UpdateZone(id, &zone)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedZone)

		case http.MethodDelete:
			err := h.zoneService.DeleteZone(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /shelves - Obtener todos los estantes o con filtros
// POST /shelves - Crear un nuevo estante
func (h *ShelfHandler) HandleShelves(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.ShelfFilter{}

			code := r.URL.Query().Get("code")
			if code != "" {
				filter.Code = code
				hasFilters = true
			}

			zoneIDStr := r.URL.Query().Get("zone_id")
			if zoneIDStr != "" {
				zoneID, err := strconv.ParseInt(zoneIDStr, 10, 64)
				if err != nil || zoneID <= 0 {
					http.Error(w, "El parámetro zone_id es inválido", http.StatusBadRequest)
					return
				}

				filter.ZoneID = &zoneID
				hasFilters = true
			}

			var shelves []*models.Shelf
			var err error

			if hasFilters {
				shelves, err = h.shelfService.GetShelvesFiltered(filter)
			} else {
				shelves, err = h.shelfService.GetAllShelves()
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shelves)

		case http.MethodPost:
			var shelf models.Shelf
			err := json.NewDecoder(r.Body).Decode(&shelf)
			if err != nil {
				http.Error(w, "Datos de estante inválidos", http.StatusBadRequest)
				return
			}

			createdShelf, err := h.shelfService.CreateShelf(&shelf)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdShelf)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /shelves/{id} - Obtener estante por ID
// PUT /shelves/{id} - Actualizar estante por ID
// DELETE /shelves/{id} - Eliminar estante por ID
func (h *ShelfHandler) HandleShelfByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/shelves/")
	if idParam == "" {
		http.Error(w, "El parámetro ID es requerido", http.StatusBadRequest)
		return
	}

	readId, err := strconv.Atoi(idParam)
	if err != nil || readId <= 0 {
		http.Error(w, "El ID es inválido", http.StatusBadRequest)
		return
	}

	id := int64(readId)

	switch r.Method {
		case http.MethodGet:
			shelf, err := h.shelfService.GetShelfByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shelf)

		case http.MethodPut:
			var shelf models.Shelf
			err := json.NewDecoder(r.Body).Decode(&shelf)
			if err != nil {
				http.Error(w, "Datos de estante inválidos", http.StatusBadRequest)
				return
			}

			updatedShelf, err := h.shelfService.UpdateShelf(id, &shelf)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedShelf)

		case http.MethodDelete:
			err := h.shelfService.DeleteShelf(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /copies - Obtener todas las copias o con filtros
// POST /copies - Crear una nueva copia
func (h *CopyHandler) HandleCopies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.CopyFilter{}

			code := r.URL.Query().Get("code")
			if code != "" {
				filter.Code = code
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

			condition := r.URL.Query().Get("condition")
			if condition != "" {
				filter.Condition = condition
				hasFilters = true
			}

			var copies []*models.Copy
			var err error

			if hasFilters {
				copies, err = h.copyService.GetCopiesFiltered(filter)
			} else {
				copies, err = h.copyService.GetAllCopies()
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(copies)

		case http.MethodPost:
			var copy models.Copy
			err := json.NewDecoder(r.Body).Decode(&copy)
			if err != nil {
				http.Error(w, "Datos de copia inválidos", http.StatusBadRequest)
				return
			}

			createdCopy, err := h.copyService.CreateCopy(&copy)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdCopy)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /copies/{id} - Obtener copia por ID
// PUT /copies/{id} - Actualizar copia por ID
// DELETE /copies/{id} - Eliminar copia por ID
func (h *CopyHandler) HandleCopyByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/copies/")
	if idParam == "" {
		http.Error(w, "El parámetro ID es requerido", http.StatusBadRequest)
		return
	}

	readId, err := strconv.Atoi(idParam)
	if err != nil || readId <= 0 {
		http.Error(w, "El ID es inválido", http.StatusBadRequest)
		return
	}

	id := int64(readId)

	switch r.Method {
		case http.MethodGet:
			copy, err := h.copyService.GetCopyByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(copy)

		case http.MethodPut:
			var copy models.Copy
			err := json.NewDecoder(r.Body).Decode(&copy)
			if err != nil {
				http.Error(w, "Datos de copia inválidos", http.StatusBadRequest)
				return
			}

			updatedCopy, err := h.copyService.UpdateCopy(id, &copy)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCopy)

		case http.MethodDelete:
			err := h.copyService.DeleteCopy(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}
