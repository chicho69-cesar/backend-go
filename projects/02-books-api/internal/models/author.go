package models

import "database/sql"

type Author struct {
	ID          int64          `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Biography   sql.NullString `json:"biography"`
	Nationality sql.NullString `json:"nationality"`
}
