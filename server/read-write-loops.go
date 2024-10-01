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
				// チャンネルがクローズされた場合、WebSocketを閉じる
				log.Println("Client send channel closed")
				client.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// メッセージを送信
			err := client.Ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// エラーが発生した場合
				log.Printf("Error sending message to client %s: %v", client.ID, err)
				return
			}
		}
	}
}

// Clientの初期化
func clientInitialize(client *Client) {
	time.Sleep(1 * time.Second)
	client.initialized = true
}
