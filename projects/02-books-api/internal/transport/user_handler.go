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

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GET /users - Obtener todos los usuarios o con filtros
// POST /users - Crear un nuevo usuario
func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var hasFilters bool
			filter := store.UserFilter{}

			code := r.URL.Query().Get("code")
			if code != "" {
				filter.Code = code
				hasFilters = true
			}

			dni := r.URL.Query().Get("dni")
			if dni != "" {
				filter.DNI = dni
				hasFilters = true
			}

			userType := r.URL.Query().Get("user_type")
			if userType != "" {
				filter.UserType = userType
				hasFilters = true
			}

			status := r.URL.Query().Get("status")
			if status != "" {
				filter.Status = status
				hasFilters = true
			}

			var users []*models.User
			var err error

			if hasFilters {
				users, err = h.userService.GetUsersFiltered(filter)
			} else {
				users, err = h.userService.GetAllUsers()
			}

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)

		case http.MethodPost:
			var user models.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				http.Error(w, "Datos de usuario inválidos", http.StatusBadRequest)
				return
			}

			createdUser, err := h.userService.CreateUser(&user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(createdUser)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}

// GET /users/{id} - Obtener usuario por ID
// PUT /users/{id} - Actualizar usuario por ID
// DELETE /users/{id} - Eliminar usuario por ID
func (h *UserHandler) HandleUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/users/")
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
			user, err := h.userService.GetUserByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)

		case http.MethodPut:
			var user models.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				http.Error(w, "Datos de usuario inválidos", http.StatusBadRequest)
				return
			}

			updatedUser, err := h.userService.UpdateUser(id, &user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedUser)

		case http.MethodDelete:
			err := h.userService.DeleteUser(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "Unavailable Method", http.StatusMethodNotAllowed)
	}
}
