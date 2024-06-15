package main

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github/pikachu0310/hackathon24spring-server/server"
	"log"
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
	log.Fatal(e.Start(":1729"))
}
