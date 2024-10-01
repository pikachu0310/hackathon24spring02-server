package server

import (
	"encoding/json"
	"github/pikachu0310/hackathon24spring-server/domain"
	"log"
)

// プレイヤー情報の送信
func broadcastPlayerUpdate(playerData domain.PlayerData) {
	playerDataJSON, err := json.Marshal(playerData)
	if err != nil {
		log.Printf("Error marshalling player data: %v", err)
		return
	}

	sendToOtherClientByPlayerID(playerDataJSON, playerData.ID)
}

// アイテム情報の送信
func broadcastItemUpdate(itemData domain.ItemData) {
	itemDataJSON, err := json.Marshal(itemData)
	if err != nil {
		log.Printf("Error marshalling item data: %v", err)
		return
	}

	sendToAllClients(itemDataJSON)
}

// 弾丸情報の送信
func broadcastBulletUpdate(bulletData domain.BulletData) {
	bulletDataJSON, err := json.Marshal(bulletData)
	if err != nil {
		log.Printf("Error marshalling bullet data: %v", err)
		return
	}

	sendToOtherClientByPlayerID(bulletDataJSON, bulletData.ShooterID)
}

// 全クライアントにデータを送信
func sendToAllClients(message []byte) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	for _, client := range clients {
		if client.initialized {
			client.send <- message
		}
	}
}

// 自分以外の全クライアントにデータを送信
func sendToOtherClientByPlayerID(message []byte, playerID string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	for _, client := range clients {
		if clientIDToPlayerID[client.ID] != playerID && client.initialized {
			client.send <- message
		}
	}
}

func sendToOtherClients(message []byte, clientID string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	for _, client := range clients {
		if client.ID != clientID && client.initialized {
			client.send <- message
		}
	}
}
