package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
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
	println(ws.Request().FormValue("number"))
	fmt.Println("New incoming connection from client: ", ws.RemoteAddr())
	s.connections[ws] = ws.Request().FormValue("number")
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
				break
			}
			fmt.Println("Read error: ", err)
			continue
		}
		if messageData.Receiver == "all" {
			s.broadcastMessage(messageData.Message)
		} else {
			for client, number := range s.connections {
				if number == messageData.Receiver {
					s.sendMessageToSpecificClient(messageData.Message, client)
				}
			}
		}
	}
}

func (s *Server) broadcastMessage(message string) {
	for ws := range s.connections {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write([]byte(message)); err != nil {
				fmt.Println("Write error: ", err)
			}
		}(ws)
	}
}

func (s *Server) sendMessageToSpecificClient(message string, receiver *websocket.Conn) {
	_, err := receiver.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending data to: "+receiver.LocalAddr().String(), err)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error while creating server: ", err)
	}
	log.Println("Server is up on port 8080")
}
