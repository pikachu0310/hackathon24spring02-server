package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	Ws          *websocket.Conn
	ID          string
	send        chan []byte
	initialized bool
}

var clients = make(map[string]*Client)
var clientMutex = &sync.Mutex{}

func NewClient(ws *websocket.Conn) *Client {
	client := &Client{
		Ws:          ws,
		ID:          uuid.New().String(),
		send:        make(chan []byte),
		initialized: false,
	}

	// WebSocketクローズハンドラの設定
	ws.SetCloseHandler(func(code int, text string) error {
		RemoveClient(client)
		return nil
	})

	// ループ処理を別ファイルの関数へ移行
	go clientReadLoop(client)
	go clientWriteLoop(client)
	go clientInitialize(client)

	return client
}

func AddClient(client *Client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	clients[client.ID] = client
	fmt.Println("New client connected: " + client.ID)
}

func RemoveClient(client *Client) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	if _, ok := clients[client.ID]; ok {
		delete(clients, client.ID)
		playerID, exists := clientIDToPlayerID[client.ID]
		if exists {
			RemovePlayer(playerID)
			RemoveClientPlayerID(client.ID)
		}
		close(client.send)
	}
}
