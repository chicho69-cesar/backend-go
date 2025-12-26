package models

import "github.com/chicho69-cesar/backend-go/books/internal/database"

type Shelf struct {
	ID          int64               `json:"id"`
	Code        string              `json:"code"` // Example: A1-01
	ZoneID      int64               `json:"zone_id"`
	Description database.NullString `json:"description"`
	LibraryID   int64               `json:"library_id"`
}
