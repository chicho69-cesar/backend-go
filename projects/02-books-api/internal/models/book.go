package models

import (
	"database/sql"
	"time"
)

type Book struct {
	ID               int64          `json:"id"`
	ISBN             string         `json:"isbn"`
	Title            string         `json:"title"`
	Subtitle         sql.NullString `json:"subtitle"`
	Edition          sql.NullString `json:"edition"`
	Language         sql.NullString `json:"language"`
	PublicationYear  sql.NullInt64  `json:"publication_year"`
	Pages            sql.NullInt64  `json:"pages"`
	Synopsis         sql.NullString `json:"synopsis"`
	PublisherID      sql.NullInt64  `json:"publisher_id"`
	ShelfID          sql.NullInt64  `json:"shelf_id"`
	Status           string         `json:"status"` // Available, Borrowed, Reserved, Maintenance
	RegistrationDate time.Time      `json:"registration_date"`
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
