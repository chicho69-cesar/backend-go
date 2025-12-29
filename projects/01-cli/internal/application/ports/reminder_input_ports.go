package ports

import (
	"context"

	"github.com/chicho69-cesar/backend-go/notes-tui/internal/domain/entities"
)

type ReminderStats struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
	Overdue   int `json:"overdue"`
	Today     int `json:"today"`
	ThisWeek  int `json:"this_week"`
}

type ReminderInputPort interface {
	CreateReminder(ctx context.Context, title, description string, dueDate entities.TimeRange, priority entities.Priority, noteID string) (*entities.Reminder, error)
	GetReminder(ctx context.Context, id string) (*entities.Reminder, error)
	GetAllReminders(ctx context.Context) ([]*entities.Reminder, error)
	GetActiveReminders(ctx context.Context) ([]*entities.Reminder, error)
	GetOverdueReminders(ctx context.Context) ([]*entities.Reminder, error)
	GetRemindersByNote(ctx context.Context, noteID string) ([]*entities.Reminder, error)
	UpdateReminder(ctx context.Context, id, title, description string, dueDate entities.TimeRange, priority entities.Priority) (*entities.Reminder, error)
	DeleteReminder(ctx context.Context, id string) error
	MarkReminderAsCompleted(ctx context.Context, id string) error
	MarkReminderAsPending(ctx context.Context, id string) error
	GetReminderStats(ctx context.Context) (ReminderStats, error)
}
