package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func Run() {
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Err listening: ", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Err Accepting con: ", err)
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read Error: %v", err)
			break
		}

		ackMsg := strings.ToUpper(strings.TrimSpace(message))
		response := fmt.Sprintf("ACK : %s\n", ackMsg)
		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Server write error: %v", err)
			break
		}
	}
}
