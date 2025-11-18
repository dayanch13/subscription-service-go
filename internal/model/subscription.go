package model

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          int       `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name" binding:"required"`
	Price       int       `json:"price" db:"price" binding:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" db:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" db:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type SubscriptionCreate struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,min=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type SubscriptionUpdate struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type CostRequest struct {
	StartPeriod string     `json:"start_period" binding:"required"`
	EndPeriod   string     `json:"end_period" binding:"required"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty"`
}

type CostResponse struct {
	TotalCost   int        `json:"total_cost"`
	Period      string     `json:"period"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty"`
}
