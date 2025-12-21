package models

import "database/sql"

type LibraryZone struct {
	ID          int64          `json:"id"`
	Code        string         `json:"code"` // Example: A1, B2, C3
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Floor       int            `json:"floor"`
}
