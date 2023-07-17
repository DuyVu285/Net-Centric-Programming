package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// User struct represents a user entity
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Temporary storage for user data
var users []User

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Hello endpoint - GET request
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}).Methods("GET")

	// Users endpoint - GET request
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}).Methods("GET")

	// Get user by ID endpoint - GET request
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		for _, user := range users {
			if user.ID == params["id"] {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(user)
				return
			}
		}
		// If the user is not found, return 404 status code
		w.WriteHeader(http.StatusNotFound)
	}).Methods("GET")

	// Add user endpoint - POST request
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		var newUser User
		_ = json.NewDecoder(r.Body).Decode(&newUser)
		users = append(users, newUser)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // Set the status code to 201
		json.NewEncoder(w).Encode(newUser)
	}).Methods("POST")

	// Update user endpoint - PUT request
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		for index, user := range users {
			if user.ID == params["id"] {
				var updatedUser User
				_ = json.NewDecoder(r.Body).Decode(&updatedUser)
				users[index] = updatedUser
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updatedUser)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}).Methods("PUT")

	// Delete user endpoint - DELETE request
	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		for index, user := range users {
			if user.ID == params["id"] {
				// Remove the user from the users slice
				users = append(users[:index], users[index+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		// If the user is not found, return 404 status code
		w.WriteHeader(http.StatusNotFound)
	}).Methods("DELETE")

	// Delete all users endpoint - DELETE request
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Delete all users logic
		users = []User{}

		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	// ...

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
