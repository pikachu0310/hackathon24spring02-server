package server

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github/pikachu0310/hackathon24spring-server/domain"
	"log"
	"sync"
	"time"
)

type Client struct {
	Ws   *websocket.Conn
	ID   string
	send chan []byte
}

var clients = make(map[*Client]bool)
var players = make(map[string]domain.PlayerData)
var mutex = &sync.Mutex{}

func NewClient(ws *websocket.Conn) *Client {
	client := &Client{
		Ws:   ws,
		ID:   generateID(),
		send: make(chan []byte),
	}

	go client.readLoop()
	go client.writeLoop()

	return client
}

func AddClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	clients[client] = true
}

func RemoveClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clients[client]; ok {
		delete(clients, client)
		close(client.send)
		delete(players, client.ID)
	}
}

func (client *Client) SendText(text string) {
	fmt.Println("[SEND] " + text)
	client.send <- []byte(text)
}

func (client *Client) readLoop() {
	defer func() {
		RemoveClient(client)
		client.Ws.Close()
	}()

	for {
		_, data, err := client.Ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		fmt.Println("[RECEIVE] " + string(data))

		var playerData domain.PlayerData
		if err := json.Unmarshal(data, &playerData); err != nil {
			log.Printf("Error unmarshalling player data: %v", err)
			continue
		}

		mutex.Lock()
		players[playerData.ID] = playerData
		mutex.Unlock()
	}
}

func (client *Client) writeLoop() {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.Ws.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func handleMessages() {
	ticker := time.NewTicker(time.Second / 40)
	defer ticker.Stop()

	for {
		<-ticker.C

		mutex.Lock()
		playerList := make([]domain.PlayerData, 0, len(players))
		for _, player := range players {
			playerList = append(playerList, player)
		}
		mutex.Unlock()

		playerDataJSON, err := json.Marshal(playerList)
		if err != nil {
			log.Printf("Error marshalling player data: %v", err)
			continue
		}

		mutex.Lock()
		for client := range clients {
			client.send <- playerDataJSON
		}
		mutex.Unlock()
	}
}

func generateID() string {
	return uuid.New().String()
}

func init() {
	go handleMessages()
}
