package models

import (
	"errors"
)

var ErrNoRecord = errors.New("No matching record found")

type ErrorResponse struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Source struct {
		Pointer string `json:"pointer"`
	} `json:"source"`
}
