package models

import "github.com/chicho69-cesar/backend-go/books/internal/database"

type Category struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Description database.NullString `json:"description"`
}
