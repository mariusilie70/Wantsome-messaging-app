package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var shutdown os.Signal = syscall.SIGUSR1

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunServer() {
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	server := &http.Server{Addr: ":8080"}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Printf("Starting server on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("error starting server: %s", err)
			stop <- shutdown
		}
	}()

	signal := <-stop
	log.Printf("Shutting down server ... ")

	m.Lock()
	for conn := range userConnections {
		conn.Close()
		delete(userConnections, conn)
	}
	m.Unlock()

	server.Shutdown(nil)
	if signal == shutdown {
		os.Exit(1)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("Home page requested")
	fmt.Fprintf(w, "Hello world from my server!")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %s", err)
		return
	}

	log.Println("New client connected")

	// Register new client connection
	m.Lock()
	userConnections[conn] = ""
	m.Unlock()

	// Handle incoming messages from the client
	go handleIncomingMessages(conn)

	// Send welcome message to the client
	welcomeMessage := models.Message{
		Sender:    "Server",
		Recipient: "",
		Content:   "Welcome to the messaging app!",
	}
	conn.WriteJSON(welcomeMessage)
}

func handleIncomingMessages(conn *websocket.Conn) {
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error reading message: %s", err)
			break
		}

		log.Printf("Received message from client: %+v", message)

		if message.Recipient == "" {
			broadcast <- message
		} else {
			sendMessageToRecipient(message)
		}
	}

	conn.Close()
	log.Println("Client disconnected")
	m.Lock()
	delete(userConnections, conn)
	m.Unlock()
}

func sendMessageToRecipient(message models.Message) {
	m.Lock()
	defer m.Unlock()

	for conn := range userConnections {
		if userConnections[conn] == message.Recipient {
			conn.WriteJSON(message)
			break
		}
	}
}

func handleMessages() {
	for {
		message := <-broadcast

		m.Lock()
		for conn := range userConnections {
			conn.WriteJSON(message)
		}
		m.Unlock()
	}
}
