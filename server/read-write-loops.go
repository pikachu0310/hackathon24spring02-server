package server

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// ClientのreadLoop
func clientReadLoop(client *Client) {
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

		handleMessage(data, client)
	}
}

// ClientのwriteLoop
func clientWriteLoop(client *Client) {
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

// Clientの初期化
func clientInitialize(client *Client) {
	time.Sleep(1 * time.Second)
	client.initialized = true
}
