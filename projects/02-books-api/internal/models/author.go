package models

import "github.com/chicho69-cesar/backend-go/books/internal/database"

type Author struct {
	ID          int64               `json:"id"`
	FirstName   string              `json:"first_name"`
	LastName    string              `json:"last_name"`
	Biography   database.NullString `json:"biography"`
	Nationality database.NullString `json:"nationality"`
	LibraryID   int64               `json:"library_id"`
}
