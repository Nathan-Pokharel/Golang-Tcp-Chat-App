package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connected to server. Enter messages to send.")

	// Start a goroutine to receive messages from the server
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error receiving message from server:", err)
				break
			}
			fmt.Print("Received message:", message)
		}
	}()

	// Read messages from the user and send them to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()

		// Send the message to the server
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message to server:", err)
			break
		}

		// If the user enters "quit", exit the client
		if strings.ToLower(message) == "quit" {
			break
		}
	}

	fmt.Println("Client disconnected.")
}

