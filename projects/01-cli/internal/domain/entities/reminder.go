package entities

import "time"

type Reminder struct {
	ID          string    `json:"id"`
	NoteID      string    `json:"note_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	IsCompleted bool      `json:"is_completed"`
	Priority    Priority  `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
	PriorityUrgent
)

func (p Priority) String() string {
	return [...]string{"Low", "Medium", "High", "Urgent"}[p]
}

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

func (r *Reminder) Update(title, description string, dueDate time.Time, priority Priority) {
	r.Title = title
	r.Description = description
	r.DueDate = dueDate
	r.Priority = priority
	r.UpdatedAt = time.Now()
}

func (r *Reminder) MarkAsCompleted() {
	r.IsCompleted = true
	r.UpdatedAt = time.Now()
}

func (r *Reminder) MarkAsPending() {
	r.IsCompleted = false
	r.UpdatedAt = time.Now()
}

func (r *Reminder) IsOverdue() bool {
	return !r.IsCompleted && time.Now().After(r.DueDate)
}

func (r *Reminder) DaysUntilDue() int {
	hoursUntil := time.Until(r.DueDate).Hours()
	return int(hoursUntil / 24)
}

func (r *Reminder) IsValid() bool {
	return r.Title != "" && !r.DueDate.IsZero()
}

func generateReminderID() string {
	return "rem-" + time.Now().Format("20060102150405") + "-" + randomString(6)
}
