package main

import (
	"fmt"
	"net"
)

func main() {
	HOST := "127.0.0.1"
	PORT := "8080"
	l, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer l.Close()

	fmt.Println("Listening on " + HOST + ":" + PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connected by", conn.RemoteAddr())

	request := make([]byte, 1024)
	message, err := conn.Read(request)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Println("Received:", string(request[:message]))
}
