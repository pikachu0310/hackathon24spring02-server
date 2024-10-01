package server

import (
	"encoding/json"
	"github/pikachu0310/hackathon24spring-server/domain"
	"log"
)

// メッセージの種類ごとに関数を切り分け
func handleMessage(data []byte, client *Client) {
	var base struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		log.Printf("Error unmarshalling base data: %v", err)
		return
	}

	switch base.Type {
	case "player":
		handlePlayerMessage(data, client)
	case "item":
		handleItemMessage(data)
	case "bullet":
		handleBulletMessage(data)
	}
}

func handlePlayerMessage(data []byte, client *Client) {
	var playerData domain.PlayerData
	if err := json.Unmarshal(data, &playerData); err != nil {
		log.Printf("Error unmarshalling player data: %v", err)
		return
	}

	// 状態管理ファイルにプレイヤーをセット
	SetPlayer(playerData)
	SetClientPlayerID(client.ID, playerData.ID)

	// プレイヤー情報を送信
	broadcastPlayerUpdate(playerData)
}

func handleItemMessage(data []byte) {
	var itemData domain.ItemData
	if err := json.Unmarshal(data, &itemData); err != nil {
		log.Printf("Error unmarshalling item data: %v", err)
		return
	}

	// 状態管理ファイルにアイテムをセット
	SetItem(itemData)

	// アイテム情報を送信
	broadcastItemUpdate(itemData)
}

func handleBulletMessage(data []byte) {
	var bulletData domain.BulletData
	if err := json.Unmarshal(data, &bulletData); err != nil {
		log.Printf("Error unmarshalling bullet data: %v", err)
		return
	}

	// 状態管理ファイルに弾丸をセット
	// SetBullet(bulletData)

	// 弾丸情報を送信
	broadcastBulletUpdate(bulletData)
}
