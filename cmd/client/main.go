package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go readMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if strings.ToLower(message) == "quit" {
			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func readMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		fmt.Printf("Received message from server: %s\n", message)
	}
}
