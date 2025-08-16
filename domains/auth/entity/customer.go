package entity

import "github.com/google/uuid"

type Customer struct {
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
	UserName    string    `json:"user_name"`
	Email       string    `json:"email"`
}
