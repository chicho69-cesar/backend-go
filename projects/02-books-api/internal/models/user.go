package models

import (
	"time"

	"github.com/chicho69-cesar/backend-go/books/internal/database"
)

type User struct {
	ID               int64               `json:"id"`
	Code             string              `json:"code"`
	DNI              string              `json:"dni"`
	FirstName        string              `json:"first_name"`
	LastName         string              `json:"last_name"`
	Email            database.NullString `json:"email"`
	Phone            database.NullString `json:"phone"`
	Address          database.NullString `json:"address"`
	UserType         string              `json:"user_type"` // Student, Teacher, Staff, External
	Status           string              `json:"status"`    // Active, Suspended, Inactive
	RegistrationDate time.Time           `json:"registration_date"`
	LibraryID        int64               `json:"library_id"`
}
