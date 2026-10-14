package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn net.Conn) {
	// In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nDisconnected from the server.")
			os.Exit(0) // Exit the client if the server shuts down
		}
		fmt.Print(msg)
	}
}

func write(conn net.Conn) {
	// Continually get input from the user and send messages to the server.
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input.")
			continue
		}
		fmt.Fprintf(conn, "%s", msg)
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	// Try to connect to the server
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		fmt.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to", *addrPtr)

	// Start asynchronously reading and displaying messages
	go read(conn)

	// Start getting and sending user messages.
	// We run this directly in main() so it blocks and keeps the program alive.
	write(conn)
}
