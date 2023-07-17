package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer conn.Close()

	fmt.Println(">> Connected to server.")
	conn.Write([]byte("Client Connected\n"))

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		fmt.Print("->: ", message)

		if strings.Contains(strings.TrimSpace(message), "Ending") {
			fmt.Println("Do you want to play again? Yes/No")
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">>: ")
			text, _ := reader.ReadString('\n')

			if strings.Contains(text, "Yes") {
				conn.Write([]byte("Yes\n"))
			} else {
				conn.Write([]byte("No\n"))
				fmt.Println("The game ends!")
				return
			}
		} else {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">>: ")
			text, _ := reader.ReadString('\n')
			fmt.Fprintf(conn, text)
			if strings.TrimSpace(text) == "STOP" {
				fmt.Println("Player quits!")
				fmt.Println("TCP client exiting...")
				return
			}
		}
	}
}
