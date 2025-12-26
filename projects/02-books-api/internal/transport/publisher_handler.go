package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
)

type PublisherHandler struct {
	publisherService *services.PublisherService
}

func NewPublisherHandler(publisherService *services.PublisherService) *PublisherHandler {
	return &PublisherHandler{
		publisherService: publisherService,
	}
}

// GET /publishers - Obtener todas las editoriales
// POST /publishers - Crear una nueva editorial
func (h *PublisherHandler) HandlePublishers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			publishers, err := h.publisherService.GetAllPublishers()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(publishers)
		
		case http.MethodPost:
			var publisher models.Publisher
			err := json.NewDecoder(r.Body).Decode(&publisher)
			if err != nil {
				http.Error(w, "Datos de editorial inválidos", http.StatusBadRequest)
				return
			}

			createdPublisher, err := h.publisherService.CreatePublisher(&publisher)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdPublisher)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}

// GET /publishers/{id} - Obtener una editorial por ID
// PUT /publishers/{id} - Actualizar una editorial por ID
// DELETE /publishers/{id} - Eliminar una editorial por ID
func (h *PublisherHandler) HandlePublisherByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/publishers/"):]
	if idParam == "" {
		http.Error(w, "El parámetro ID es requerido", http.StatusBadRequest)
		return
	}

	readId, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "El ID es inválido", http.StatusBadRequest)
		return
	}

	id := int64(readId)

	switch r.Method {
		case http.MethodGet:
			publisher, err := h.publisherService.GetPublisherByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(publisher)

		case http.MethodPut:
			var publisher models.Publisher
			err := json.NewDecoder(r.Body).Decode(&publisher)
			if err != nil {
				http.Error(w, "Datos de editorial inválidos", http.StatusBadRequest)
				return
			}

			updatedPublisher, err := h.publisherService.UpdatePublisher(id, &publisher)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedPublisher)

		case http.MethodDelete:
			err := h.publisherService.DeletePublisher(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}
