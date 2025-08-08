package entity

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	CustomerID     uuid.UUID  `json:"customer_id" db:"customer_id"`
	Email          string     `json:"email" db:"email"`
	FirstName      string     `json:"first_name" db:"first_name"`
	LastName       string     `json:"last_name" db:"last_name"`
	Phone          *string    `json:"phone" db:"phone"`
	Address        *string    `json:"address" db:"address"`
	City           *string    `json:"city" db:"city"`
	PostalCode     *string    `json:"postal_code" db:"postal_code"`
	OpenIDSub      string     `json:"openid_sub" db:"openid_sub"`
	SetupCompleted bool       `json:"setup_completed" db:"setup_completed"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`
}

type CustomerStats struct {
	TotalCustomers    int64   `json:"total_customers"`
	CompletedSetup    int64   `json:"completed_setup"`
	PendingSetup      int64   `json:"pending_setup"`
	AverageOrderValue float64 `json:"average_order_value"`
	TotalOrdersPlaced int64   `json:"total_orders_placed"`
}

type CustomerProfile struct {
	Customer
	OrderCount       int64      `json:"order_count"`
	TotalSpent       float64    `json:"total_spent"`
	AverageOrderSize float64    `json:"average_order_size"`
	LastOrderDate    *time.Time `json:"last_order_date,omitempty"`
}
