package models

import (
	"github.com/google/uuid"
)

type CreateTicketOptionRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Allocation  int    `json:"allocation"`
		} `json:"attributes"`
	} `json:"data"`
}

type TicketOptionResponse struct {
	Data struct {
		Type       string    `json:"type"`
		ID         uuid.UUID `json:"id"`
		Attributes struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Allocation  int    `json:"allocation"`
		} `json:"attributes"`
	} `json:"data"`
}

type CreatePurchaseRequest struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Quantity int `json:"quantity"`
		} `json:"attributes"`
		Relationships struct {
			TicketOption struct {
				Data struct {
					Type string    `json:"type"`
					ID   uuid.UUID `json:"id"`
				} `json:"data"`
			} `json:"ticket_option"`
			User struct {
				Data struct {
					Type string    `json:"type"`
					ID   uuid.UUID `json:"id"`
				} `json:"data"`
			} `json:"user"`
		} `json:"relationships"`
	} `json:"data"`
}

type PurchaseResponse struct {
	Data struct {
		ID         uuid.UUID `json:"id"`
		Type       string    `json:"type"`
		Attributes struct {
			Quantity int `json:"quantity"`
		} `json:"attributes"`
		Relationships struct {
			TicketOption struct {
				Data struct {
					Type string    `json:"type"`
					ID   uuid.UUID `json:"id"`
				} `json:"data"`
			} `json:"ticket_option"`
			User struct {
				Data struct {
					Type string    `json:"type"`
					ID   uuid.UUID `json:"id"`
				} `json:"data"`
			} `json:"user"`
		} `json:"relationships"`
	} `json:"data"`
}
