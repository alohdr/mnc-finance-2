package models

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
)
