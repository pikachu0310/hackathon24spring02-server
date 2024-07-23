package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github/pikachu0310/hackathon24spring-server/server"
	"log"
	"net"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnectionRequest(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return err
	}

	client := server.NewClient(ws)
	server.AddClient(client)

	return nil
}

func main() {
	e := echo.New()
	e.GET("/ws", HandleConnectionRequest)

	// IPv4とIPv6の両方でリッスンする
	listener, err := net.Listen("tcp", "[::]:1729")
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	log.Fatal(http.Serve(listener, e))
}
