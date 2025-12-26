package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chicho69-cesar/backend-go/books/internal/models"
	"github.com/chicho69-cesar/backend-go/books/internal/store"
	"github.com/chicho69-cesar/backend-go/books/internal/validations"
)

type PublisherService struct {
	publisherStore store.IPublisherStore
}

func NewPublisherService(publisherStore store.IPublisherStore) *PublisherService {
	return &PublisherService{
		publisherStore: publisherStore,
	}
}

func (s *PublisherService) GetAllPublishers(libraryID int64) ([]*models.Publisher, error) {
	return s.publisherStore.GetAll(libraryID)
}

func (s *PublisherService) GetPublisherByID(libraryID, id int64) (*models.Publisher, error) {
	if id <= 0 {
		return nil, errors.New("El id de la editorial es invalido")
	}

	publisher, err := s.publisherStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la editorial con ID %d: %w", id, err)
	}

	return publisher, nil
}

func (s *PublisherService) CreatePublisher(libraryID int64, publisher *models.Publisher) (*models.Publisher, error) {
	if err := validations.ValidatePublisher(publisher); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	publisher.LibraryID = libraryID
	publisher.Name = strings.TrimSpace(publisher.Name)

	if publisher.Country.Valid {
		publisher.Country.String = strings.TrimSpace(publisher.Country.String)
	}

	createdPublisher, err := s.publisherStore.Create(libraryID, publisher)
	if err != nil {
		return nil, fmt.Errorf("Error al crear la editorial: %w", err)
	}

	return createdPublisher, nil
}

func (s *PublisherService) UpdatePublisher(libraryID, id int64, publisher *models.Publisher) (*models.Publisher, error) {
	if id <= 0 {
		return nil, errors.New("El id de la editorial es invalido")
	}

	existingPublisher, err := s.publisherStore.GetByID(libraryID, id)
	if err != nil {
		return nil, fmt.Errorf("Error al obtener la editorial con ID %d: %w", id, err)
	}

	if existingPublisher == nil {
		return nil, fmt.Errorf("La editorial con ID %d no existe", id)
	}

	if err := validations.ValidatePublisher(publisher); err != nil {
		return nil, fmt.Errorf("Validación fallida: %w", err)
	}

	publisher.Name = strings.TrimSpace(publisher.Name)

	if publisher.Country.Valid {
		publisher.Country.String = strings.TrimSpace(publisher.Country.String)
	}

	updatedPublisher, err := s.publisherStore.Update(libraryID, id, publisher)
	if err != nil {
		return nil, fmt.Errorf("Error al actualizar la editorial con ID %d: %w", id, err)
	}

	return updatedPublisher, nil
}

func (s *PublisherService) DeletePublisher(libraryID, id int64) error {
	if id <= 0 {
		return errors.New("El id de la editorial es invalido")
	}

	existingPublisher, err := s.publisherStore.GetByID(libraryID, id)
	if err != nil {
		return fmt.Errorf("Error al obtener la editorial con ID %d: %w", id, err)
	}

	if existingPublisher == nil {
		return fmt.Errorf("La editorial con ID %d no existe", id)
	}

	if err := s.publisherStore.Delete(libraryID, id); err != nil {
		return fmt.Errorf("Error al eliminar la editorial con ID %d: %w", id, err)
	}

	return nil
}
