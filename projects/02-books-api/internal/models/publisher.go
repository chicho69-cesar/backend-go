package models

import "github.com/chicho69-cesar/backend-go/books/internal/database"

type Publisher struct {
	ID      int64               `json:"id"`
	Name    string              `json:"name"`
	Country database.NullString `json:"country"`
}
