package services

import (
	"fmt"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
)

type ConfigurationService struct {
	configStore store.IConfigStore
}

func NewConfigurationService(configStore store.IConfigStore) *ConfigurationService {
	return &ConfigurationService{
		configStore: configStore,
	}
}

func (s *ConfigurationService) GetConfiguration() (*models.Configuration, error) {
	config, err := s.configStore.GetCurrent()
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la configuración: %w", err)
	}

	return config, nil
}

func (s *ConfigurationService) UpdateConfiguration(updates map[string]any) (*models.Configuration, error) {
	currentConfig, err := s.configStore.GetCurrent()
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la configuración actual: %w", err)
	}

	if val, ok := updates["student_loan_days"]; ok {
		if intVal, ok := val.(float64); ok {
			currentConfig.StudentLoanDays = int(intVal)
		}
	}

	if val, ok := updates["teacher_loan_days"]; ok {
		if intVal, ok := val.(float64); ok {
			currentConfig.TeacherLoanDays = int(intVal)
		}
	}

	if val, ok := updates["max_renewals"]; ok {
		if intVal, ok := val.(float64); ok {
			currentConfig.MaxRenewals = int(intVal)
		}
	}

	if val, ok := updates["max_books_per_loan"]; ok {
		if intVal, ok := val.(float64); ok {
			currentConfig.MaxBooksPerLoan = int(intVal)
		}
	}

	if val, ok := updates["fine_per_day"]; ok {
		if floatVal, ok := val.(float64); ok {
			currentConfig.FinePerDay = floatVal
		}
	}

	if val, ok := updates["reservation_days"]; ok {
		if intVal, ok := val.(float64); ok {
			currentConfig.ReservationDays = int(intVal)
		}
	}

	if val, ok := updates["grace_days"]; ok {
		if intVal, ok := val.(float64); ok {
			currentConfig.GraceDays = int(intVal)
		}
	}

	updatedConfig, err := s.configStore.Update(currentConfig)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar la configuración: %w", err)
	}

	return updatedConfig, nil
}
