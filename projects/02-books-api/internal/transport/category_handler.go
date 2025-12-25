package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			categories, err := h.categoryService.GetAllCategories()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		
		case http.MethodPost:
			var category models.Category
			err := json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				http.Error(w, "Datos de categoría inválidos", http.StatusBadRequest)
				return
			}

			createdCategory, err := h.categoryService.CreateCategory(&category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdCategory)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Path[len("/categories/"):]
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
			category, err := h.categoryService.GetCategoryByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)

		case http.MethodPut:
			var category models.Category
			err := json.NewDecoder(r.Body).Decode(&category)
			if err != nil {
				http.Error(w, "Datos de categoría inválidos", http.StatusBadRequest)
				return
			}

			updatedCategory, err := h.categoryService.UpdateCategory(id, &category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)

		case http.MethodDelete:
			err := h.categoryService.DeleteCategory(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}
