package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fixedTime, _ := time.Parse(time.RFC3339, "2024-07-19T12:00:00Z")
	tickets := []Ticket{
		{EventName: "Concert A", Venue: "Venue A", EventDate: fixedTime, Price: 50.0, Num: 100},
		{EventName: "Concert B", Venue: "Venue B", EventDate: fixedTime, Price: 75.0, Num: 200},
	}

	for i, ticket := range tickets {
		ticketAsBytes, err := json.Marshal(ticket)
		if err != nil {
			return fmt.Errorf("failed to marshal ticket %d: %s", i, err.Error())
		}

		err = ctx.GetStub().PutState(fmt.Sprintf("TICKET%d", i), ticketAsBytes)
		if err != nil {
			return fmt.Errorf("failed to put ticket %d: %s", i, err.Error())
		}
	}

	return nil
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

	ticketAsBytes, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %s", err.Error())
	}
	return ctx.GetStub().PutState(id, ticketAsBytes)
}

// QueryTicket allows anyone to query a ticket by ID
func (s *SmartContract) QueryTicket(ctx contractapi.TransactionContextInterface, id string) (*Ticket, error) {
	log.Printf("QueryTicket called with id: %s", id)

	ticketAsBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		log.Printf("Failed to read from world state: %s", err.Error())
		return nil, fmt.Errorf("failed to read from world state: %s", err.Error())
	}

	if ticketAsBytes == nil {
		log.Printf("Ticket with id %s does not exist", id)
		return nil, fmt.Errorf("%s does not exist", id)
	}

	ticket := new(Ticket)
	err = json.Unmarshal(ticketAsBytes, ticket)
	if err != nil {
		log.Printf("Failed to unmarshal ticket: %s", err.Error())
		return nil, fmt.Errorf("failed to unmarshal ticket: %s", err.Error())
	}

	log.Printf("Ticket found: %+v", ticket)
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
		err = json.Unmarshal(queryResponse.Value, ticket)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal ticket: %s", err.Error())
		}

		queryResult := QueryResult{Key: queryResponse.Key, Record: ticket}
		results = append(results, queryResult)
	}

	return results, nil
}

// CreateOrder allows a customer to create a new order (purchase a ticket)
func (s *SmartContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderID string, userID string, ticketID string, quantity int, totalPrice float64, orderDate string) error {
	// Parse orderDate from string to time.Time
	orderTime, err := time.Parse(time.RFC3339, orderDate)
	if err != nil {
		return fmt.Errorf("failed to parse order date: %s", err.Error())
	}

	// Query the ticket to check availability and price
	ticket, err := s.QueryTicket(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to query ticket: %s", err.Error())
	}

	// Check if there are enough tickets available
	if ticket.Num < uint(quantity) {
		return fmt.Errorf("not enough tickets available")
	}

	// Verify the total price
	calculatedTotalPrice := float64(quantity) * ticket.Price
	if totalPrice != calculatedTotalPrice {
		return fmt.Errorf("total price mismatch: expected %.2f, got %.2f", calculatedTotalPrice, totalPrice)
	}

	// Update the ticket quantity
	ticket.Num -= uint(quantity)
	ticketAsBytes, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %s", err.Error())
	}
	err = ctx.GetStub().PutState(ticketID, ticketAsBytes)
	if err != nil {
		return fmt.Errorf("failed to update ticket state: %s", err.Error())
	}

	// Create the Order struct
	order := Order{
		UserID:     userID,
		TicketID:   ticketID,
		Quantity:   quantity,
		TotalPrice: totalPrice,
		OrderDate:  orderTime,
	}

	// Marshal the Order struct to JSON
	orderAsBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %s", err.Error())
	}

	// Put the order state in the ledger
	return ctx.GetStub().PutState(orderID, orderAsBytes)
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("创建链码失败: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("启动链码失败: %s", err.Error())
	}
}
