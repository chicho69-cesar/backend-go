package models

import (
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
)

type Book struct {
	ID               int64               `json:"id"`
	ISBN             string              `json:"isbn"`
	Title            string              `json:"title"`
	Subtitle         database.NullString `json:"subtitle"`
	Edition          database.NullString `json:"edition"`
	Language         database.NullString `json:"language"`
	PublicationYear  database.NullInt64  `json:"publication_year"`
	Pages            database.NullInt64  `json:"pages"`
	Synopsis         database.NullString `json:"synopsis"`
	PublisherID      database.NullInt64  `json:"publisher_id"`
	ShelfID          database.NullInt64  `json:"shelf_id"`
	Status           string              `json:"status"` // Available, Borrowed, Reserved, Maintenance
	RegistrationDate time.Time           `json:"registration_date"`
}

type BookAuthor struct {
	BookID   int64 `json:"book_id"`
	AuthorID int64 `json:"author_id"`
	Position int   `json:"position"` // Author order (1st author, 2nd author, etc.)
}

type BookCategory struct {
	BookID     int64 `json:"book_id"`
	CategoryID int64 `json:"category_id"`
}
