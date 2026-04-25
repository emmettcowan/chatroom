package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Clients struct {
	mu    sync.RWMutex
	conns map[net.Conn]net.Conn
}

func (clients *Clients) addClient(conn net.Conn) {
	clients.mu.Lock()
	clients.conns[conn] = conn
	clients.mu.Unlock()
}

func (clients *Clients) removeClient(conn net.Conn) {
	clients.mu.Lock()
	delete(clients.conns, conn)
	clients.mu.Unlock()
}

func (clients *Clients) broadCast(message string) {
	clients.mu.RLock()
	for _, conn := range clients.conns {
		fmt.Print(message)
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Printf("Server write error: %v", err)
			continue
		}
	}
	clients.mu.RUnlock()
}

func Run() {
	clients := Clients{conns: map[net.Conn]net.Conn{}}
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
		clients.addClient(conn)

		log.Printf("Clients connected %d", len(clients.conns))

		go handleConnection(conn, &clients)

	}
}

func handleConnection(conn net.Conn, clients *Clients) {
	defer conn.Close()
	defer clients.removeClient(conn)

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read Error: %v", err)
			break
		}

		ackMsg := strings.TrimSpace(message)
		response := fmt.Sprintf("%s\n", ackMsg)
		fmt.Println(response)
		clients.broadCast(response)
	}
}
