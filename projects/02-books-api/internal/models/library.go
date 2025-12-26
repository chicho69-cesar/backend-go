package models

type Library struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Website  string `json:"website"`
	Username string `json:"username"`
	Password string `json:"password"`
}
