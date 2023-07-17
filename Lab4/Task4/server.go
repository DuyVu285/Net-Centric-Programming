package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
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
		return
	}
	if message == 0 {
		return
	}
	fmt.Println("Received:", string(request[:message]))

	extract := strings.Split(string(request[:message]), " ")
	filePath := extract[1]

	var content []byte
	var contentType string

	switch {
	case strings.HasSuffix(filePath, ".html"):
		content, err = ioutil.ReadFile("." + filePath)
		if err != nil {
			fmt.Println(err)
			content = []byte("404 Not Found")
			contentType = "text/plain"
			break
		}
		contentType = "text/html"
	case strings.HasSuffix(filePath, ".jpg"):
		content, err = ioutil.ReadFile("." + filePath)
		if err != nil {
			fmt.Println(err)
			content = []byte("404 Not Found")
			contentType = "text/plain"
			break
		}
		contentType = "image/jpeg"
	case strings.HasSuffix(filePath, ".mp3"):
		content, err = ioutil.ReadFile("." + filePath)
		if err != nil {
			fmt.Println(err)
			content = []byte("404 Not Found")
			contentType = "text/plain"
			break
		}
		contentType = "audio/mpeg"
	case strings.HasSuffix(filePath, ".wav"):
		content, err = ioutil.ReadFile("." + filePath)
		if err != nil {
			fmt.Println(err)
			content = []byte("404 Not Found")
			contentType = "text/plain"
			break
		}
		contentType = "audio/wav"
	case strings.HasSuffix(filePath, ".mid"):
		content, err = ioutil.ReadFile("." + filePath)
		if err != nil {
			fmt.Println(err)
			content = []byte("404 Not Found")
			contentType = "text/plain"
			break
		}
		contentType = "audio/midi"
	default:
		content = []byte("404 Not Found")
		contentType = "text/plain"
	}

	var response bytes.Buffer
	response.WriteString("HTTP/1.1 200 OK\r\n")
	response.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
	response.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(content)))
	response.WriteString("\r\n")
	response.Write(content)

	conn.Write(response.Bytes())
}
