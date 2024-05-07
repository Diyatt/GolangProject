package handlers

import (
	"fmt"
	"net/http"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a protected resource!")
}
