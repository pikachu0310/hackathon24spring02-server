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
	Ws          *websocket.Conn
	ID          string
	send        chan []byte
	initialized bool
}

var clients = make(map[string]*Client)
var players = make(map[string]domain.PlayerData)
var items = make(map[string]domain.ItemData)
var bullets = make(map[string]domain.BulletData)
var clientIDToPlayerID = make(map[string]string)
var mutex = &sync.Mutex{}

func NewClient(ws *websocket.Conn) *Client {
	client := &Client{
		Ws:          ws,
		ID:          generateID(),
		send:        make(chan []byte),
		initialized: false,
	}

	// Closeハンドラの設定
	ws.SetCloseHandler(func(code int, text string) error {
		log.Printf("WebSocket closed with code %d and message %s", code, text)
		RemoveClient(client)
		return nil
	})

	go client.readLoop()
	go client.writeLoop()
	go client.initialize()

	return client
}

func AddClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	clients[client.ID] = client
	fmt.Println("New client connected: " + client.ID)
}

func RemoveClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := clients[client.ID]; ok {
		delete(clients, client.ID)
		close(client.send)

		// clientIDToPlayerIDからplayerIDを取得し、playersから削除
		if playerID, ok := clientIDToPlayerID[client.ID]; ok {
			delete(players, playerID)
			delete(clientIDToPlayerID, client.ID)
		}

		fmt.Println("Client disconnected: " + client.ID)
	}
}

func (client *Client) SendText(text string) {
	//fmt.Println("[SEND] " + text)
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
		//fmt.Println("[RECEIVE] " + string(data))

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
			clientIDToPlayerID[client.ID] = playerData.ID // clientIDとplayerIDを紐づける
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

func (client *Client) initialize() {
	time.Sleep(3 * time.Second)
	client.initialized = true
}

func handleMessages() {
	ticker := time.NewTicker(time.Second / 30)
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

		for _, client := range clients {
			if client.initialized {
				client.send <- playerDataJSON
				client.send <- itemDataJSON
				client.send <- bulletDataJSON
			}
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
