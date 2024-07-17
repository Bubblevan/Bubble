package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID     uint
	TicketID   uint
	Quantity   int
	TotalPrice float64
}
