package models

import (
	"github.com/google/uuid"
)

type TicketOption struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Allocation  int       `json:"allocation"`
}

type Purchase struct {
	ID           uuid.UUID `json:"id"`
	TicketOption uuid.UUID `json:"ticket_option_id"`
	User         uuid.UUID `json:"user_id"`
	Quantity     int       `json:"quantity"`
}

type User struct {
	ID uuid.UUID `json:"id"`
}

type Ticket struct {
	ID       uuid.UUID `json:"id"`
	Purchase uuid.UUID `json:"purchase_id"`
}
