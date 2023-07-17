package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// User struct represents a user entity
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	// Create a new user
	newUser := User{
		ID:       "1",
		Username: "john_doe",
		Email:    "john@example.com",
	}

	newUser2 := User{
		ID:       "2",
		Username: "john_doe2",
		Email:    "john2@example.com",
	}
	// Send a POST request to create the user
	createUser(newUser)
	createUser(newUser2)
	//deleteUser(newUser.ID)
	//deleteAllUsers()

	// Update the user
	/*updatedUser := User{
		ID:       "1",
		Username: "johndoe",
		Email:    "johndoe@example.com",
	}*/

	// Send a PUT request to update the user
	//updateUser(updatedUser)

	// Send a DELETE request to delete the user
	//deleteUser(updatedUser.ID)

	//deleteAllUsers()
}

func createUser(user User) {
	url := "http://localhost:8080/users"

	// Convert the user struct to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	// Send the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("User created successfully")
	} else {
		fmt.Println("Failed to create user. Status:", resp.Status)
	}
}

func updateUser(user User) {
	url := "http://localhost:8080/users/" + user.ID

	// Convert the user struct to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	// Create a PUT request
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the PUT request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("User updated successfully")
	} else {
		fmt.Println("Failed to update user. Status:", resp.Status)
	}
}

func deleteUser(userID string) {
	url := "http://localhost:8080/users/" + userID

	// Create a DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Send the DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("User deleted successfully")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("User not found")
	} else {
		fmt.Println("Failed to delete user. Status:", resp.Status)
	}
}

func deleteAllUsers() {
	url := "http://localhost:8080/users"

	// Create a DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Send the DELETE request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("All users deleted successfully")
	} else {
		fmt.Println("Failed to delete all users. Status:", resp.Status)
	}
}
