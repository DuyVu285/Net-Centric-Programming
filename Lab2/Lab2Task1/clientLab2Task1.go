package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	var getPrefix bool
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to Server")

	// Authenticate with the server
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	fmt.Fprintf(conn, username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	fmt.Fprintf(conn, password)

	response, _ := bufio.NewReader(conn).ReadString('\n')
	if response == "Authentication failed\n" {
		fmt.Println(response)
		return
	} else {
		re := regexp.MustCompile(`\d+`)
		prefix := re.FindString(response)
		fmt.Println("Prefix: ", prefix)
		getPrefix = true
		if getPrefix {
			for {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter message: ")
				message, _ := reader.ReadString('\n')
				message_prefix := prefix + "_" + message
				request := fmt.Sprintf(message_prefix)
				fmt.Fprintf(conn, request)
				if strings.Contains(message, "STOP") {
					fmt.Println("TCP client exiting...")
					return
				}

				response, _ := bufio.NewReader(conn).ReadString('\n')
				fmt.Print("Server response: " + response)
				if strings.Contains(response, "STOP") {
					fmt.Println("TCP server exited...")
					return
				}
			}
		}
	}
}
