package main

import (
"bufio"
"fmt"
"net"
"os"
)

type Client struct {
conn   net.Conn
writer *bufio.Writer
}

func main() {
// Start the TCP server
listener, err := net.Listen("tcp", ":8080")
if err != nil {
    fmt.Println("Error starting the server:", err)
    os.Exit(1)
}
defer listener.Close()
fmt.Println("Server started. Listening on :8080")

// Channel to broadcast messages to all connected clients
messages := make(chan string)

// Channel to register new clients
register := make(chan Client)

// Channel to unregister disconnected clients
unregister := make(chan Client)

// Goroutine to handle broadcasting messages to all clients
go func() {
    clients := make(map[net.Conn]struct{})

    for {
        select {
        case msg := <-messages:
            // Send the message to all connected clients
            for client := range clients {
                _, err := client.Write([]byte(msg))
                if err != nil {
                    fmt.Println("Error broadcasting message to client:", err)
                }
            }

        case client := <-register:
            // Register a new client
            clients[client.conn] = struct{}{}
            fmt.Println("New client connected:", client.conn.RemoteAddr())

        case client := <-unregister:
            // Unregister a client
            delete(clients, client.conn)
            fmt.Println("Client disconnected:", client.conn.RemoteAddr())
        }
    }
}()

// Listen for incoming connections
for {
    conn, err := listener.Accept()
    if err != nil {
        fmt.Println("Error accepting connection:", err)
        continue
    }

    // Create a new client and register it
    client := Client{conn: conn, writer: bufio.NewWriter(conn)}
    register <- client

    // Goroutine to handle receiving messages from the client
    go func(client Client) {
        reader := bufio.NewReader(client.conn)

        for {
            // Read a message from the client
            message, err := reader.ReadString('\n')
            if err != nil {
                break
            }

            // Broadcast the message to all clients
            messages <- message
        }

        // Unregister the client when the connection is closed
        unregister <- client
        client.conn.Close()
    }(client)
}
}

