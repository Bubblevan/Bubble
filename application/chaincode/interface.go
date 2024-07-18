package main

import (
	"time"
)

type Ticket struct {
	EventName string    `json:"eventName"`
	Venue     string    `json:"venue"`
	EventDate time.Time `json:"eventDate"`
	Price     float64   `json:"price"`
	Num       uint      `json:"num"`
}

type Order struct {
	UserID     string  `json:"userId"`
	TicketID   string  `json:"ticketId"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
}

type QueryResult struct {
	Key    string `json:"key"`
	Record *Ticket `json:"record"`
}
