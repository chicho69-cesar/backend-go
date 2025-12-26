package transport

import (
	"encoding/json"
	"net/http"

	"github.com/chicho69-cesar/backend-go/books/internal/middleware"
	"github.com/chicho69-cesar/backend-go/books/internal/services"
)

type ConfigurationHandler struct {
	configService *services.ConfigurationService
}

func NewConfigurationHandler(configService *services.ConfigurationService) *ConfigurationHandler {
	return &ConfigurationHandler{
		configService: configService,
	}
}

// GET /configuration - Obtener la configuración actual
// PATCH /configuration - Actualizar la configuración
func (h *ConfigurationHandler) HandleConfiguration(w http.ResponseWriter, r *http.Request) {
	libraryID, err := middleware.GetLibraryID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
		case http.MethodGet:
			config, err := h.configService.GetConfiguration(libraryID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(config)

		case http.MethodPatch:
			var updates map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&updates)
			if err != nil {
				http.Error(w, "Datos de configuración inválidos", http.StatusBadRequest)
				return
			}

			if len(updates) == 0 {
				http.Error(w, "No se proporcionaron campos para actualizar", http.StatusBadRequest)
				return
			}

			updatedConfig, err := h.configService.UpdateConfiguration(libraryID, updates)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedConfig)

		default:
			http.Error(w, "Unavailable method", http.StatusMethodNotAllowed)
	}
}
