package models

import "database/sql"

type Shelf struct {
	ID          int64          `json:"id"`
	Code        string         `json:"code"` // Example: A1-01
	ZoneID      int64          `json:"zone_id"`
	Description sql.NullString `json:"description"`
}
