package repositories

import (
	"context"

	"github.com/chicho69-cesar/backend-go/notes-tui/internal/domain/entities"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder *entities.Reminder) error
	FindByID(ctx context.Context, id string) (*entities.Reminder, error)
	FindAll(ctx context.Context) ([]*entities.Reminder, error)
	FindByNoteID(ctx context.Context, noteID string) ([]*entities.Reminder, error)
	FindActive(ctx context.Context) ([]*entities.Reminder, error)
	FindCompleted(ctx context.Context) ([]*entities.Reminder, error)
	FindOverdue(ctx context.Context) ([]*entities.Reminder, error)
	FindByDueDateRange(ctx context.Context, start, end entities.TimeRange) ([]*entities.Reminder, error)
	Update(ctx context.Context, reminder *entities.Reminder) error
	Delete(ctx context.Context, id string) error
	MarkAsCompleted(ctx context.Context, id string) error
	MarkAsPending(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
	CountOverdue(ctx context.Context) (int, error)
}
