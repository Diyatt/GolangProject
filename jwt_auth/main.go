package main

import (
	"fmt"
	"jwt_auth/handlers"
	"jwt_auth/middleware"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", handlers.SignUpHandler)
	http.HandleFunc("/signin", handlers.SignInHandler)
	http.HandleFunc("/protected", middleware.AuthMiddleware(handlers.ProtectedHandler))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
