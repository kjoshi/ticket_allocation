package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"encoding/json"

	"tickets/models"
)

func (app *application) GetTicketOption(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	ticketOption, err := app.db.GetTicketOption(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "database_error", "Database error", err.Error(), "")
		}
		return
	}

	resp := models.TicketOptionResponse{}
	resp.Data.Type = "ticket_options"
	resp.Data.ID = ticketOption.ID
	resp.Data.Attributes.Name = ticketOption.Name
	resp.Data.Attributes.Description = ticketOption.Description
	resp.Data.Attributes.Allocation = ticketOption.Allocation

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (app *application) CreateTicketOption(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req models.CreateTicketOptionRequest

	err := dec.Decode(&req)
	if err != nil {
		// FIXME: Returning err.Error() gives away too much detail about how the API is implemented
		sendBadRequestResponse(w, "invalid_request", "Invalid request body", err.Error(), "/data")
		return
	}

	if req.Data.Type != "ticket_options" {
		sendBadRequestResponse(w, "invalid_type", "Invalid data type", "Expected 'ticket_options'", "/data/type")
		return
	}

	if req.Data.Attributes.Name == "" {
		sendBadRequestResponse(w, "invalid_name", "Invalid attribute: name", "Name must not be empty", "/data/attributes/name")
		return
	}

	if req.Data.Attributes.Allocation < 1 {
		sendBadRequestResponse(w, "invalid_allocation", "Invalid attribute: allocation", "Allocation must be greater than zero", "/data/attributes/allocation")
		return
	}

	ticketOption, err :=
		app.db.CreateTicketOption(req.Data.Attributes.Name, req.Data.Attributes.Description, req.Data.Attributes.Allocation)

	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "database_error", "Database error", err.Error(), "")
		return
	}

	resp := models.TicketOptionResponse{}
	resp.Data.Type = "ticket_options"
	resp.Data.ID = ticketOption.ID
	resp.Data.Attributes.Name = ticketOption.Name
	resp.Data.Attributes.Description = ticketOption.Description
	resp.Data.Attributes.Allocation = ticketOption.Allocation

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (app *application) PurchaseTickets(w http.ResponseWriter, r *http.Request) {
	// TODO: Create a helper to decode the JSON body
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req models.CreatePurchaseRequest
	err := dec.Decode(&req)
	if err != nil {
		sendBadRequestResponse(w, "invalid_request", "Invalid request body", err.Error(), "/data")
		return
	}

	if req.Data.Type != "purchases" {
		sendBadRequestResponse(w, "invalid_type", "Invalid data type", "Expected 'purchases'", "/data/type")
		return
	}

	if req.Data.Attributes.Quantity < 1 {
		sendBadRequestResponse(w, "invalid_quantity", "Invalid quantity", "Quantity must be greater than zero", "/data/attributes/quantity")
		return
	}

	_, err = app.db.GetTicketOption(req.Data.Relationships.TicketOption.Data.ID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			sendBadRequestResponse(w, "invalid_ticket_option", "Invalid ticket option", "Ticket option not found", "/data/relationships/ticket_option/data/id")
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "database_error", "Database error", err.Error(), "")
		}
		return
	}

	//TODO: Check that the user exists

	purchase, err := app.db.CreatePurchase(req.Data.Relationships.TicketOption.Data.ID, req.Data.Attributes.Quantity, req.Data.Relationships.User.Data.ID)

	if err != nil {
		if errors.Is(err, models.ErrInvalidQuantity) {
			sendBadRequestResponse(w, "invalid_purchase_quantity", "Unable to purhase quantity provided", "Unable to reserve given quantity of ticket options", "/data/attributes/quantity")
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "database_error", "Database error", err.Error(), "")
		}
		return
	}
	resp := models.PurchaseResponse{}
	resp.Data.Type = "purchases"
	resp.Data.ID = purchase.ID
	resp.Data.Attributes.Quantity = purchase.Quantity
	resp.Data.Relationships.TicketOption.Data.Type = "ticket_options"
	resp.Data.Relationships.TicketOption.Data.ID = purchase.TicketOption
	resp.Data.Relationships.User.Data.Type = "users"
	resp.Data.Relationships.User.Data.ID = purchase.User

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}
