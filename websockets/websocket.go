package mywebsockets

import (
	"fmt"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func (s *Server) broadcast(b []byte, room string) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if room == ws.Request().URL.Query()["id"][0] {
				_, err := ws.Write(b)
				if err != nil {
					fmt.Println("write error:", err)
				}
			}

		}(ws)
	}
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}
