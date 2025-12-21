package models

import "time"

type Reservation struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	BookID          int64     `json:"book_id"`
	ReservationDate time.Time `json:"reservation_date"`
	ExpirationDate  time.Time `json:"expiration_date"`
	Status          string    `json:"status"` // Pending, Active, Cancelled, Expired
	Priority        int       `json:"priority"`
	Notified        bool      `json:"notified"`
}
