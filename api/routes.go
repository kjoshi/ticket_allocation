package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ticket_options/{id}", app.GetTicketOption)
	mux.HandleFunc("POST /ticket_options", app.CreateTicketOption)

	return mux
}
