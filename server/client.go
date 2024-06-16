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
var items = make(map[string]domain.ItemData)
var bullets = make(map[string]domain.BulletData)
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

		var base struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(data, &base); err != nil {
			log.Printf("Error unmarshalling base data: %v", err)
			continue
		}

		switch base.Type {
		case "player":
			var playerData domain.PlayerData
			if err := json.Unmarshal(data, &playerData); err != nil {
				log.Printf("Error unmarshalling player data: %v", err)
				continue
			}
			mutex.Lock()
			players[playerData.ID] = playerData
			mutex.Unlock()
		case "item":
			var itemData domain.ItemData
			if err := json.Unmarshal(data, &itemData); err != nil {
				log.Printf("Error unmarshalling item data: %v", err)
				continue
			}
			mutex.Lock()
			items[itemData.ID] = itemData
			mutex.Unlock()
		case "bullet":
			var bulletData domain.BulletData
			if err := json.Unmarshal(data, &bulletData); err != nil {
				log.Printf("Error unmarshalling bullet data: %v", err)
				continue
			}
			mutex.Lock()
			bullets[bulletData.ID] = bulletData
			mutex.Unlock()
		}
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
	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()

	for {
		<-ticker.C

		mutex.Lock()
		playerList := make([]domain.PlayerData, 0, len(players))
		for _, player := range players {
			playerList = append(playerList, player)
		}
		playerDataJSON, err := json.Marshal(playerList)
		if err != nil {
			log.Printf("Error marshalling player data: %v", err)
			mutex.Unlock()
			continue
		}

		itemList := make([]domain.ItemData, 0, len(items))
		for _, item := range items {
			itemList = append(itemList, item)
		}
		itemDataJSON, err := json.Marshal(itemList)
		if err != nil {
			log.Printf("Error marshalling item data: %v", err)
			mutex.Unlock()
			continue
		}

		bulletList := make([]domain.BulletData, 0, len(bullets))
		for _, bullet := range bullets {
			bulletList = append(bulletList, bullet)
		}
		bulletDataJSON, err := json.Marshal(bulletList)
		if err != nil {
			log.Printf("Error marshalling bullet data: %v", err)
			mutex.Unlock()
			continue
		}

		// reset bullets
		bullets = make(map[string]domain.BulletData)

		for client := range clients {
			client.send <- playerDataJSON
			client.send <- itemDataJSON
			client.send <- bulletDataJSON
		}
		mutex.Unlock()
	}
}

func generateID() string {
	return uuid.New().String()
}

func init() {
	go handleMessages()
	go generateItems()
}
