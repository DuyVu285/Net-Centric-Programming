package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("udp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Client Connected!")
	fmt.Print("Enter your username: ")
	username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Println("Start Chatting!")
	// Register with the server
	fmt.Fprintf(conn, "%s", username)

	go receiveMessages(conn)

	// Read from stdin and send messages to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "@all ") || strings.HasPrefix(text, "@") {
			fmt.Fprintf(conn, "%s\n", text)
		} else {
			fmt.Println("Invalid format. Use @all or @<username> to chat.")
		}
		if text == "/quit" {
			fmt.Fprintf(conn, "%s\n", text)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func receiveMessages(conn net.Conn) {
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		message := string(buf[:n])
		fmt.Print(message)
	}
}
