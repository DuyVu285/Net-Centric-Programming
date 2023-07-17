package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func numberGuessing(conn net.Conn, random_num int, inputnumber int) {
	fmt.Println("Client guessed: ", inputnumber)
	if inputnumber > random_num {
		Larger := "The guessed number is larger!" + "\n"
		conn.Write([]byte(Larger))
	} else if inputnumber < random_num {
		Smaller := "The guessed number is smaller!" + "\n"
		conn.Write([]byte(Smaller))
	} else if inputnumber == random_num {
		Correct := "The guessed number is correct! Ending..." + "\n"
		conn.Write([]byte(Correct))
		fmt.Println("The game ends!")
		return
	}
}

func generateRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	random_num := 1 + rand.Intn(100-1+1)
	fmt.Printf("The random number is: %d\n", random_num)
	return random_num
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	random_number := generateRandomNumber()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print("-> ", string(netData))

		switch strings.TrimSpace(string(netData)) {
		case "Client Connected":
			Confirm := "Start Guessing Number\n"
			conn.Write([]byte(Confirm))
		case "Yes":
			Confirm := "Start Guessing Number\n"
			conn.Write([]byte(Confirm))
			random_number = generateRandomNumber()
		case "No":
			fmt.Println("Player quits!")
			fmt.Println("Exiting TCP server!")
			return
		default:
			if inputnumber, err := strconv.Atoi(strings.TrimSpace(string(netData))); err == nil {
				numberGuessing(conn, random_number, inputnumber)
			} else {
				Invalid := "Invalid number!\n"
				conn.Write([]byte(Invalid))
			}
		}

		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Player quits!")
			fmt.Println("Exiting TCP server!")
			return
		}
	}
}
func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Println("TCP server started.")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		go handleConnection(conn)
	}

}
