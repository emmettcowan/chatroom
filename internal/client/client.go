package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Run() {
	conn, err := net.Dial("tcp", ":8090")
	if err != nil {
		log.Fatal("Error connecting: ", err)
	}
	defer conn.Close()

	go readFromUser(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			if err != os.ErrClosed {
				log.Printf("Error reading stdin: %v", err)
			}
			break
		}

		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			break
		}
	}
}

func readFromUser(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nDisconnected from server: %v\n", err)
			return
		}
		fmt.Printf("\rServer: %s>: ", response)
	}
}
