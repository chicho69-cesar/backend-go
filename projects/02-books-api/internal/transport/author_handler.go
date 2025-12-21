package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
)

type AuthorHandler struct {
	authorService *services.AuthorService
}

func NewAuthorHandler(authorService *services.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		authorService: authorService,
	}
}

func (h *AuthorHandler) HandleAuthors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			authors, err := h.authorService.GetAllAuthors()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(authors)

		case http.MethodPost:
			var author models.Author
			err := json.NewDecoder(r.Body).Decode(&author)
			if err != nil {
				http.Error(w, "Datos de autor inválidos", http.StatusBadRequest)
				return
			}

			createdAuthor, err := h.authorService.CreateAuthor(&author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdAuthor)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}

func (h *AuthorHandler) HandleAuthorByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
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
			author, err := h.authorService.GetAuthorByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(author)

		case http.MethodPut:
			var author models.Author
			err := json.NewDecoder(r.Body).Decode(&author)
			if err != nil {
				http.Error(w, "Datos de autor inválidos", http.StatusBadRequest)
				return
			}

			updatedAuthor, err := h.authorService.UpdateAuthor(id, &author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application-json")
			json.NewEncoder(w).Encode(updatedAuthor)

		case http.MethodDelete:
			err := h.authorService.DeleteAuthor(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
			
		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}
