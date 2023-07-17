package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// User struct
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Email       string    `json:"email"`
	Avatar      string    `json:"avatar"`
	PhoneNumber string    `json:"phonenumber"`
	DateOfBirth time.Time `json:"dateofbirth"`
	Address     string    `json:"address"`
}

var db *sql.DB

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Database connection
	var err error
	db, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/user_management?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Hello endpoint - GET request
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}).Methods("GET")

	// Users endpoint - GET request
	router.HandleFunc("/users", getUsers).Methods("GET")

	// Get user by Endpoint - GET request
	router.HandleFunc("/users/{id}", getUserByID).Methods("GET")

	// Add user endpoint - POST request
	router.HandleFunc("/users", addUser).Methods("POST")

	// Update user endpoint - PUT request
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")

	// Delete user endpoint - DELETE request
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Delete all users endpoint - DELETE request
	router.HandleFunc("/users", deleteAllUsers).Methods("DELETE")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all users
// Get users with filtering and sorting
func getUsers(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	firstname := r.URL.Query().Get("firstname")
	lastname := r.URL.Query().Get("lastname")
	sortBy := r.URL.Query().Get("sort")
	orderBy := "id" // Default sorting attribute

	switch sortBy {
	case "username":
		orderBy = "username"
	case "firstname":
		orderBy = "firstname"
	case "lastname":
		orderBy = "lastname"
	case "email":
		orderBy = "email"
	case "dateofbirth":
		orderBy = "dateofbirth"
	case "address":
		orderBy = "address"
	}

	// Build the SQL query with filtering and sorting options
	query := "SELECT * FROM users WHERE 1=1"
	params := []interface{}{}

	if username != "" {
		query += " AND username = ?"
		params = append(params, username)
	}
	if firstname != "" {
		query += " AND firstname = ?"
		params = append(params, firstname)
	}
	if lastname != "" {
		query += " AND lastname = ?"
		params = append(params, lastname)
	}

	query += " ORDER BY " + orderBy

	rows, err := db.Query(query, params...)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Avatar, &user.PhoneNumber, &user.DateOfBirth, &user.Address)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Get user by ID
func getUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["username"]

	var user User
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Avatar, &user.PhoneNumber, &user.DateOfBirth, &user.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Add user
func isUserEmpty(user User) bool {
	// Check if any required field is empty
	return user.Username == "" || user.Email == "" || user.Firstname == "" || user.Lastname == ""
}

func generateRandomUserInfo(user *User) {
	url := "https://random-data-api.com/api/v2/users"

	// Send GET request to random-data-api.com
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Decode the response into a RandomUser struct
	var randomUser struct {
		ID          int       `json:"id"`
		Username    string    `json:"username"`
		Firstname   string    `json:"first_name"`
		Lastname    string    `json:"last_name"`
		Email       string    `json:"email"`
		Avatar      string    `json:"avatar"`
		PhoneNumber string    `json:"phonenumber"`
		DateOfBirth time.Time `json:"dateofbirth"`
		Address     string    `json:"address"`
	}
	err = json.NewDecoder(resp.Body).Decode(&randomUser)
	if err != nil {
		log.Println(err)
		return
	}

	// Assign the random user information to the User object
	user.ID = randomUser.ID
	user.Username = randomUser.Username
	user.Firstname = randomUser.Firstname
	user.Lastname = randomUser.Lastname
	user.Email = randomUser.Email
	user.Avatar = randomUser.Avatar
	user.PhoneNumber = randomUser.PhoneNumber
	user.DateOfBirth = randomUser.DateOfBirth
	user.Address = randomUser.Address
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isUserEmpty(newUser) {
		// Generate random user information
		generateRandomUserInfo(&newUser)
	}

	print(newUser.Username)

	_, err = db.Exec("INSERT INTO users (id, username, firstname, lastname, email, avatar, phonenumber, dateofbirth, address) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newUser.ID, newUser.Username, newUser.Firstname, newUser.Lastname, newUser.Email, newUser.Avatar, newUser.PhoneNumber, newUser.DateOfBirth, newUser.Address)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET username = ?, firstname = ?, lastname = ?, email = ?, avatar = ?, phonenumber = ?, dateofbirth = ?, address = ? WHERE id = ?",
		updatedUser.Username, updatedUser.Firstname, updatedUser.Lastname, updatedUser.Email, updatedUser.Avatar, updatedUser.PhoneNumber, updatedUser.DateOfBirth, updatedUser.Address, id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// Delete user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete all users
func deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
