package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type User struct {
	IPAddress *net.UDPAddr
	Username  string
}

func main() {
	var users []User
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Server listening on", addr)

	for {
		buf := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		command := string(buf[:n])
		parts := strings.Split(command, " ")
		username := command
		if strings.HasPrefix(command, "@all ") || strings.HasPrefix(command, "@") {
			username = parts[0][1:]
		}

		switch {
		case strings.HasPrefix(command, "@all "):
			message := fmt.Sprintf("[%s]: %s", findUserByAddress(addr, users), strings.Join(parts[1:], " "))
			broadcastMessage(conn, message, users)
		case strings.HasPrefix(command, "@"):
			recipient := findUserByUsername(username, users)
			if recipient != nil {
				message := fmt.Sprintf("[%s]: %s", findUserByAddress(addr, users), strings.Join(parts[1:], " "))
				sendMessage(conn, message, recipient.IPAddress)
			} else {
				sendErrorMessage(conn, addr)
			}
		case strings.HasPrefix(command, "/quit"):
			fmt.Println("User logout:", findUserByAddress(addr, users))
			users = removeUser(addr, users)
		default:
			fmt.Println("Registering user...")
			if !contains(addr, users) {
				newUser := User{
					IPAddress: addr,
					Username:  username,
				}
				users = append(users, newUser)
				fmt.Println("User registered:", newUser.Username)
			}
		}
	}
}

func removeUser(addr *net.UDPAddr, users []User) []User {
	newUsers := []User{}
	for _, user := range users {
		if addr.IP.Equal(user.IPAddress.IP) && addr.Port == user.IPAddress.Port {
			continue
		}
		newUsers = append(newUsers, user)
	}
	return newUsers
}
func findUserByUsername(username string, users []User) *User {
	for _, user := range users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

func findUserByAddress(addr *net.UDPAddr, users []User) string {
	for _, user := range users {
		if addr.IP.Equal(user.IPAddress.IP) && addr.Port == user.IPAddress.Port {
			return user.Username
		}
	}
	return "Error"
}

func broadcastMessage(conn *net.UDPConn, message string, users []User) {
	fmt.Println("Broadcasting message...")
	fmt.Print(message)
	for _, user := range users {
		if _, err := conn.WriteToUDP([]byte(message), user.IPAddress); err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Message Broadcasted!")
}

func sendMessage(conn *net.UDPConn, message string, addr *net.UDPAddr) {
	fmt.Println("Sending message to client...")
	fmt.Print(message)
	if _, err := conn.WriteToUDP([]byte(message), addr); err != nil {
		log.Println(err)
	}
	fmt.Println("Message Sent!")
}

func sendErrorMessage(conn *net.UDPConn, addr *net.UDPAddr) {
	fmt.Println("Sending message...")
	message := "[server] :There is currently no user with that name. \n"
	if _, err := conn.WriteToUDP([]byte(message), addr); err != nil {
		log.Println(err)
	}
	fmt.Println("Message Sent!")
}

func contains(addr *net.UDPAddr, users []User) bool {
	for _, user := range users {
		if addr.String() == user.IPAddress.String() {
			return true
		}
	}
	return false
}
