package ports

import (
	"context"

	"github.com/chicho69-cesar/backend-go/notes-tui/internal/domain/entities"
)

type NoteStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Archived int `json:"archived"`
	WithTags int `json:"with_tags"`
	ThisWeek int `json:"this_week"`
}

type NoteInputPort interface {
	CreateNote(ctx context.Context, title, content string, tags []string) (*entities.Note, error)
	GetNote(ctx context.Context, id string) (*entities.Note, error)
	GetAllNotes(ctx context.Context) ([]*entities.Note, error)
	GetActiveNotes(ctx context.Context) ([]*entities.Note, error)
	GetArchivedNotes(ctx context.Context) ([]*entities.Note, error)
	UpdateNote(ctx context.Context, id, title, content string, tags []string) (*entities.Note, error)
	DeleteNote(ctx context.Context, id string) error
	ArchiveNote(ctx context.Context, id string) error
	UnarchiveNote(ctx context.Context, id string) error
	SearchNotes(ctx context.Context, query string) ([]*entities.Note, error)
	GetStats(ctx context.Context) (NoteStats, error)
}
