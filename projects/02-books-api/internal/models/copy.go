package models

import "github.com/chicho69-cesar/backend-go/books/internal/database"

type Copy struct {
	ID              int64                `json:"id"`
	Code            string               `json:"code"` // Barcode
	BookID          int64                `json:"book_id"`
	Status          string               `json:"status"`    // Available, Borrowed, Reserved, Damaged, Lost
	Condition       string               `json:"condition"` // New, Good, Fair, Poor
	AcquisitionDate database.NullTime    `json:"acquisition_date"`
	PurchasePrice   database.NullFloat64 `json:"purchase_price"`
	Notes           database.NullString  `json:"notes"`
}
