package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint      `json:"userID"`
	TicketID   uint      `json:"ticketID"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"totalPrice"`
	OrderDate  time.Time `json:"orderDate"`
}
