package main

import (
	"fmt"
	"io"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client: ", ws.RemoteAddr())

	// TODO: For security reasons, make sure to have a mutex to avoid race conditions
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:n]

		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Write error: ", err)
			}
		}(ws)
	}
}

func (s *Server) liveFeed(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client to live feed: ", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("%d:%d:%d\n", time.Now().Hour(), time.Now().Minute(), time.Now().Second())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 1)
	}
}
