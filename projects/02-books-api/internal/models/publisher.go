package models

import "database/sql"

type Publisher struct {
	ID      int64          `json:"id"`
	Name    string         `json:"name"`
	Country sql.NullString `json:"country"`
}
