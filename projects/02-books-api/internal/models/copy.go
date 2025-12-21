package models

import "database/sql"

type Copy struct {
	ID              int64           `json:"id"`
	Code            string          `json:"code"` // Barcode
	BookID          int64           `json:"book_id"`
	Status          string          `json:"status"`    // Available, Borrowed, Reserved, Damaged, Lost
	Condition       string          `json:"condition"` // New, Good, Fair, Poor
	AcquisitionDate sql.NullTime    `json:"acquisition_date"`
	PurchasePrice   sql.NullFloat64 `json:"purchase_price"`
	Notes           sql.NullString  `json:"notes"`
}
