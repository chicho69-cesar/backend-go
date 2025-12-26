package models

import "github.com/chicho69-cesar/backend-go/books/internal/database"

type LibraryZone struct {
	ID          int64               `json:"id"`
	Code        string              `json:"code"` // Example: A1, B2, C3
	Name        string              `json:"name"`
	Description database.NullString `json:"description"`
	Floor       int                 `json:"floor"`
	LibraryID   int64               `json:"library_id"`
}
