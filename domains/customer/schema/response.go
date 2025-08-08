package schema

import (
	"time"

	"github.com/google/uuid"
)

type CustomerResponse struct {
	CustomerID     uuid.UUID  `json:"customer_id"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	Phone          *string    `json:"phone,omitempty"`
	Address        *string    `json:"address,omitempty"`
	City           *string    `json:"city,omitempty"`
	PostalCode     *string    `json:"postal_code,omitempty"`
	SetupCompleted bool       `json:"setup_completed"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type CustomerProfileResponse struct {
	CustomerResponse
	OrderCount       int64      `json:"order_count"`
	TotalSpent       float64    `json:"total_spent"`
	AverageOrderSize float64    `json:"average_order_size"`
	LastOrderDate    *time.Time `json:"last_order_date,omitempty"`
}

type CustomerListResponse struct {
	Customers []*CustomerResponse `json:"customers"`
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	Limit     int                 `json:"limit"`
	HasNext   bool                `json:"has_next"`
}

type CustomerSearchResponse struct {
	Customers []*CustomerResponse `json:"customers"`
	Total     int64               `json:"total"`
	Page      int                 `json:"page"`
	Limit     int                 `json:"limit"`
	HasNext   bool                `json:"has_next"`
	Query     string              `json:"query,omitempty"`
}

type CustomerStatsResponse struct {
	TotalCustomers      int64   `json:"total_customers"`
	CompletedSetup      int64   `json:"completed_setup"`
	PendingSetup        int64   `json:"pending_setup"`
	SetupCompletionRate float64 `json:"setup_completion_rate"`
	AverageOrderValue   float64 `json:"average_order_value"`
	TotalOrdersPlaced   int64   `json:"total_orders_placed"`
}

type CustomerDashboardResponse struct {
	Profile         *CustomerProfileResponse `json:"profile"`
	RecentOrders    interface{}              `json:"recent_orders"`   // Will be populated from order domain
	Recommendations interface{}              `json:"recommendations"` // For future features
}
