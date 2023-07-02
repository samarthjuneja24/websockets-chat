package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	connections map[*websocket.Conn]string
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]string),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client: ", ws.RemoteAddr())
	s.connections[ws] = ws.Request().FormValue("number")
	go s.pingPong(ws)
	s.readLoop(ws)
}

type MessageData struct {
	Message  string `json:"message"`
	Receiver string `json:"receiver"`
}

func (s *Server) readLoop(ws *websocket.Conn) {
	for {
		var messageData MessageData
		err := websocket.JSON.Receive(ws, &messageData)
		if err != nil {
			//If connection from client breaks itself
			if err == io.EOF {
				delete(s.connections, ws)
				break
			}
			fmt.Println("Read error: ", err)
			continue
		}
		if messageData.Receiver == "all" {
			fmt.Println("Broadcasting to all users")
			senderNumber := s.connections[ws]
			s.broadcastMessage(messageData.Message, senderNumber)
		} else {
			fmt.Println("Sending to individual user:", messageData.Receiver)
			for client, number := range s.connections {
				if number == messageData.Receiver {
					fmt.Println("Sending message to this client:", messageData.Receiver)
					s.sendMessageToSpecificClient(messageData.Message, client)
				}
			}
		}
	}
}

func (s *Server) broadcastMessage(message string, senderNumber string) {
	for ws := range s.connections {
		if s.connections[ws] != senderNumber {
			go func(ws *websocket.Conn) {
				if _, err := ws.Write([]byte(message)); err != nil {
					fmt.Println("Write error: ", err)
				}
			}(ws)
		}
	}
}

func (s *Server) sendMessageToSpecificClient(message string, receiver *websocket.Conn) {
	_, err := receiver.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending data to: "+receiver.LocalAddr().String(), err)
	}
}

func (s *Server) pingPong(conn *websocket.Conn) {
	pingInterval := time.Second * 60
	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()
	for range ticker.C {
		err := websocket.Message.Send(conn, "ping")
		if err != nil {
			log.Println("Error while ping-pong to client: ", err)
			delete(s.connections, conn)
			break
		} else {
			log.Println("Connection successful for client:", conn.LocalAddr().String())
		}
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	println(port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("Error while creating server: ", err)
	}
	log.Println("Server is up on port " + port)
}
