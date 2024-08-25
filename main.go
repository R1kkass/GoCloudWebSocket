package main

import (
	"log"
	"net/http"

	"mypackages/db"
	mywebsockets "mypackages/websockets"

	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

type ConnBool map[*websocket.Conn]bool

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	db.ConnectDatabase()

	server := mywebsockets.NewServer()

	config := &websocket.Config{
		Origin: nil,
	}
	
	http.Handle("/ws", websocket.Server{Handler: server.HandleChatWs, Config: *config})

	http.ListenAndServe(":8000", nil)
}
