package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")

	stdinReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		// 1. Prompt and read input from user
		fmt.Print("Enter message: ")
		text, _ := stdinReader.ReadString('\n')

		// 2. Send over connection
		fmt.Fprintf(conn, text)

		// 3. Read acknowledgement back from server
		reply, _ := serverReader.ReadString('\n')
		fmt.Print("Reply: " + reply)
	}
}
