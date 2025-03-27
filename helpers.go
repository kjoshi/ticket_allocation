package main

import (
	"encoding/json"
	"net/http"
	"tickets/models"
)

func sendErrorResponse(w http.ResponseWriter, status int, code, title, detail, pointer string) {
	errResp := models.ErrorResponse{
		Status: status,
		Code:   code,
		Title:  title,
		Detail: detail,
	}
	errResp.Source.Pointer = pointer

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errResp)
}

func sendBadRequestResponse(w http.ResponseWriter, code, title, detail, pointer string) {
	sendErrorResponse(w, http.StatusBadRequest, code, title, detail, pointer)
}
