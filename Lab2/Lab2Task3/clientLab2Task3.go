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
		fmt.Println(err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	fmt.Print(message)

	for {
		message, _ := reader.ReadString('\n')
		fmt.Print(message)

		if strings.Contains(message, "Game over") {
			break
		}

		fmt.Print("Enter a letter: ")
		text, _ := bufio.NewReader(os.Stdin).ReadString('\n')

		conn.Write([]byte(text))
	}
}
