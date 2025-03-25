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

func (app *application) CreateTicketOption(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req models.CreateTicketOptionRequest

	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Data.Type != "ticket_options" {
		http.Error(w, "Invalid value for data.type", http.StatusBadRequest)
		return
	}

	if req.Data.Attributes.Name == "" {
		http.Error(w, "Invalid value for data.attributes.name", http.StatusBadRequest)
		return
	}

	if req.Data.Attributes.Allocation < 1 {
		http.Error(w, "Invalid value for data.attributes.allocation", http.StatusBadRequest)
		return
	}

	ticketOption, err := app.db.CreateTicketOption(req.Data.Attributes.Name, req.Data.Attributes.Description, req.Data.Attributes.Allocation)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
