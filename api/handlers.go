package main

import (
	"errors"
	"github.com/google/uuid"
	"net/http"

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
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
