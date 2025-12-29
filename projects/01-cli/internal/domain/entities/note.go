package entities

import "time"

type Note struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsArchived bool      `json:"is_archived"`
}

func NewNote(title, content string, tags []string) *Note {
	now := time.Now()

	return &Note{
		ID:         generateID(),
		Title:      title,
		Content:    content,
		Tags:       tags,
		CreatedAt:  now,
		UpdatedAt:  now,
		IsArchived: false,
	}
}

func (n *Note) Update(title, content string, tags []string) {
	n.Title = title
	n.Content = content
	n.Tags = tags
	n.UpdatedAt = time.Now()
}

func (n *Note) Archive() {
	n.IsArchived = true
	n.UpdatedAt = time.Now()
}

func (n *Note) Unarchive() {
	n.IsArchived = false
	n.UpdatedAt = time.Now()
}

func (n *Note) IsValid() bool {
	return n.Title != "" && n.Content != ""
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)

	for i := range b {
		b[i] = letters[time.Now().Nanosecond()%len(letters)]
	}

	return string(b)
}
