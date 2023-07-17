package main

import (
	"bytes"
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

	var response bytes.Buffer
	response.WriteString("HTTP/1.1 200 OK\r\n")
	response.WriteString("Content-Type: text/html\r\n")
	response.WriteString("Content-Length: 53\r\n")
	response.WriteString("\r\n<html>Welcome to Very Simple Web Server</html>")

	conn.Write(response.Bytes())
}
