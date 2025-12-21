package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID               int64          `json:"id"`
	Code             string         `json:"code"`
	DNI              string         `json:"dni"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Email            sql.NullString `json:"email"`
	Phone            sql.NullString `json:"phone"`
	Address          sql.NullString `json:"address"`
	UserType         string         `json:"user_type"` // Student, Teacher, Staff, External
	Status           string         `json:"status"`    // Active, Suspended, Inactive
	RegistrationDate time.Time      `json:"registration_date"`
}
