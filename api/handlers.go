package main

import (
	"fmt"
	"net/http"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
