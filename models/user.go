package models

import "time"

type (
	User struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		PIN         string `json:"pin"`
	}

	Login struct {
		PhoneNumber string `json:"phone_number"`
		PIN         string `json:"pin"`
	}

	Profile struct {
		UserID     string    `json:"user_id"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		Address    string    `json:"address"`
		UpdateDate time.Time `json:"update_date"`
	}
)
