package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	FullName  string   `json:"full_name"`
	Emails    []string `json:"emails"`
	Addresses []string `json:"addresses"`
}

func loadUsers(fileName string) ([]User, error) {
	var users []User
	inFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	decoder := json.NewDecoder(inFile)
	err = decoder.Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func encryptPassword(password string) string {
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func decryptPassword(password string) string {
	rawDecodedText, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		panic(err)
	}
	return string(rawDecodedText)
}

func authenticateUser(username, password string, users []User) bool {
	for _, user := range users {
		fmt.Println("Authorizing...")
		if user.Username == username && decryptPassword(user.Password) == password {
			fmt.Println("Authorized")
			return true
		}
	}
	return false
}

func generatePrefix() int {
	return rand.Intn(1000)
}

func handleConnection(conn net.Conn, users []User) {
	defer conn.Close()

	// Authenticate the user
	var username, password string
	var getPrefix bool
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	username = strings.TrimSpace(name)

	pass, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	password = strings.TrimSpace(pass)

	if !authenticateUser(username, password, users) {
		fmt.Fprintln(conn, "Authentication failed")
		return
	} else {

		prefix := generatePrefix()
		fmt.Fprintf(conn, "Your prefix key is: %d\n", prefix)
		getPrefix = true
		if getPrefix {
			for {
				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Print("Message from Client: ", message)

				if strings.Contains(message, "STOP") {
					fmt.Println("TCP client exited")
					return
				}

				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter message: ")
				response, _ := reader.ReadString('\n')
				prefix_response := strconv.Itoa(prefix) + "_" + response
				request := fmt.Sprintf(prefix_response)
				fmt.Fprintf(conn, request)
				if strings.Contains(response, "STOP") {
					fmt.Println("TCP server exiting...")
					return
				}
			}
		}
	}
}

func main() {
	filePath := "user.json"
	users, err := loadUsers(filePath)
	if err != nil {
		fmt.Println("Failed to load users:", err)
		return
	}

	// Listen for incoming connections
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Connected")
		go handleConnection(conn, users)
	}
}
