package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		msg, _ := reader.ReadString('\n')

		// If the client disconnects or sends empty text, break out
		if msg == "" {
			break
		}

		fmt.Print("Received: " + msg)

		// Send an acknowledgement back over the same connection
		reply := "Server ACK: " + msg
		fmt.Fprintf(conn, reply)
	}

	conn.Close()
}

func main() {
	ln, _ := net.Listen("tcp", ":8080")

	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}
