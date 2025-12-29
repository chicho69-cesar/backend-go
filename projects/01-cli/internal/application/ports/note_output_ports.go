package ports

import (
	"context"

	"github.com/chicho69-cesar/backend-go/notes-tui/internal/domain/entities"
)

type NoteOutputPort interface {
	Create(ctx context.Context, note *entities.Note) error
	FindByID(ctx context.Context, id string) (*entities.Note, error)
	FindAll(ctx context.Context) ([]*entities.Note, error)
	FindByTags(ctx context.Context, tags []string) ([]*entities.Note, error)
	FindArchived(ctx context.Context) ([]*entities.Note, error)
	FindActive(ctx context.Context) ([]*entities.Note, error)
	Update(ctx context.Context, note *entities.Note) error
	Delete(ctx context.Context, id string) error
	Archive(ctx context.Context, id string) error
	Unarchive(ctx context.Context, id string) error
	Search(ctx context.Context, query string) ([]*entities.Note, error)
	Count(ctx context.Context) (int, error)
}
