package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// CreateTicket allows the ticket issuer to create a new ticket
func (s *SmartContract) CreateTicket(ctx contractapi.TransactionContextInterface, id string, eventName string, venue string, eventDate string, price float64, num uint) error {
	eventTime, err := time.Parse(time.RFC3339, eventDate)
	if err != nil {
		return fmt.Errorf("failed to parse event date: %s", err.Error())
	}

	ticket := Ticket{
		EventName: eventName,
		Venue:     venue,
		EventDate: eventTime,
		Price:     price,
		Num:       num,
	}

	ticketAsBytes, _ := json.Marshal(ticket)
	return ctx.GetStub().PutState(id, ticketAsBytes)
}

// QueryTicket allows anyone to query a ticket by ID
func (s *SmartContract) QueryTicket(ctx contractapi.TransactionContextInterface, id string) (*Ticket, error) {
	ticketAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %s", err.Error())
	}

	if ticketAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	ticket := new(Ticket)
	_ = json.Unmarshal(ticketAsBytes, ticket)

	return ticket, nil
}

// QueryAllTickets allows anyone to query all tickets
func (s *SmartContract) QueryAllTickets(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		ticket := new(Ticket)
		_ = json.Unmarshal(queryResponse.Value, ticket)

		queryResult := QueryResult{Key: queryResponse.Key, Record: ticket}
		results = append(results, queryResult)
	}

	return results, nil
}

// CreateOrder allows a customer to create a new order (purchase a ticket)
func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderID string, userID string, ticketID string, quantity int) error {
	ticket, err := s.QueryTicket(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket.Num < uint(quantity) {
		return fmt.Errorf("not enough tickets available")
	}

	totalPrice := float64(quantity) * ticket.Price
	order := Order{
		UserID:     userID,
		TicketID:   ticketID,
		Quantity:   quantity,
		TotalPrice: totalPrice,
	}

	ticket.Num -= uint(quantity)
	ticketAsBytes, _ := json.Marshal(ticket)
	if err := ctx.GetStub().PutState(ticketID, ticketAsBytes); err != nil {
		return fmt.Errorf("failed to update ticket state: %s", err.Error())
	}

	orderAsBytes, _ := json.Marshal(order)
	return ctx.GetStub().PutState(orderID, orderAsBytes)
}
