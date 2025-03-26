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
	return
}
