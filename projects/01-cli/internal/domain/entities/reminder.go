// internal/domain/entities/reminder.go
package entities

import (
	"time"
)

// Reminder representa un recordatorio en el dominio
type Reminder struct {
	ID          string    `json:"id"`
	NoteID      string    `json:"note_id"` // Referencia a la nota asociada (opcional)
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	IsCompleted bool      `json:"is_completed"`
	Priority    Priority  `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Priority define la prioridad de un recordatorio
type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

// String convierte Priority a string
func (p Priority) String() string {
	return [...]string{"Low", "Medium", "High", "Urgent"}[p]
}

// NewReminder crea una nueva instancia de Reminder
func NewReminder(title, description string, dueDate time.Time, priority Priority, noteID string) *Reminder {
	now := time.Now()
	return &Reminder{
		ID:          generateReminderID(),
		NoteID:      noteID,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		IsCompleted: false,
		Priority:    priority,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update modifica el recordatorio
func (r *Reminder) Update(title, description string, dueDate time.Time, priority Priority) {
	r.Title = title
	r.Description = description
	r.DueDate = dueDate
	r.Priority = priority
	r.UpdatedAt = time.Now()
}

// MarkAsCompleted marca el recordatorio como completado
func (r *Reminder) MarkAsCompleted() {
	r.IsCompleted = true
	r.UpdatedAt = time.Now()
}

// MarkAsPending marca el recordatorio como pendiente
func (r *Reminder) MarkAsPending() {
	r.IsCompleted = false
	r.UpdatedAt = time.Now()
}

// IsOverdue verifica si el recordatorio está vencido
func (r *Reminder) IsOverdue() bool {
	return !r.IsCompleted && time.Now().After(r.DueDate)
}

// DaysUntilDue días hasta la fecha límite (negativo si ya pasó)
func (r *Reminder) DaysUntilDue() int {
	hoursUntil := r.DueDate.Sub(time.Now()).Hours()
	return int(hoursUntil / 24)
}

// IsValid valida que el recordatorio tenga datos mínimos
func (r *Reminder) IsValid() bool {
	return r.Title != "" && !r.DueDate.IsZero()
}

// generateReminderID genera un ID único para el recordatorio
func generateReminderID() string {
	return "rem-" + time.Now().Format("20060102150405") + "-" + randomString(6)
}
