package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept()
		if err != nil {
			handleError(err)
			continue
		}
		// Push the newly accepted connection into the conns channel
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	reader := bufio.NewReader(client)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			// If client disconnects, stop the loop and close the connection
			client.Close()
			return
		}

		// Tidy up each message and add it to the messages channel,
		// recording which client it came from.
		msgs <- Message{
			sender:  clientid,
			message: msg,
		}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	// Create a Listener for TCP connections on the port given above.
	ln, err := net.Listen("tcp", *portPtr)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server listening on port", *portPtr)

	// Create a channel for connections
	conns := make(chan net.Conn)
	// Create a channel for messages
	msgs := make(chan Message)
	// Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	// Start accepting connections
	go acceptConns(ln, conns)

	// We need a variable to keep track of the current client ID to assign
	currentID := 0

	for {
		select {
		case conn := <-conns:
			// Deal with a new connection
			// - assign a client ID
			// - add the client to the clients map
			clients[currentID] = conn
			fmt.Printf("Client %d connected from %s\n", currentID, conn.RemoteAddr())

			// - start to asynchronously handle messages from this client
			go handleClient(conn, currentID, msgs)

			// Increment ID for the next connecting client
			currentID++

		case msg := <-msgs:
			// Deal with a new message
			// Send the message to all clients that aren't the sender
			for id, conn := range clients {
				if id != msg.sender { // Skip the sender
					fmt.Fprintf(conn, "Client %d: %s", msg.sender, msg.message)
				}
			}
		}
	}
}
