package models

type UserRegister struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	PIN         string `json:"PIN"`
}
