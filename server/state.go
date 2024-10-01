package server

import (
	"github/pikachu0310/hackathon24spring-server/domain"
	"sync"
)

// データごとに別々のmutexを使用
var (
	players            = make(map[string]domain.PlayerData)
	items              = make(map[string]domain.ItemData)
	bullets            = make(map[string]domain.BulletData)
	clientIDToPlayerID = make(map[string]string)

	playerMutex       = &sync.Mutex{} // プレイヤー用のmutex
	itemMutex         = &sync.Mutex{} // アイテム用のmutex
	bulletMutex       = &sync.Mutex{} // 弾丸用のmutex
	clientPlayerMutex = &sync.Mutex{} // clientIDとplayerID用のmutex
)

// プレイヤー関連
func SetPlayer(player domain.PlayerData) {
	playerMutex.Lock()
	defer playerMutex.Unlock()
	players[player.ID] = player
}

func GetPlayer(playerID string) (domain.PlayerData, bool) {
	playerMutex.Lock()
	defer playerMutex.Unlock()
	player, exists := players[playerID]
	return player, exists
}

func GetAllPlayers() []domain.PlayerData {
	playerMutex.Lock()
	defer playerMutex.Unlock()
	playerList := make([]domain.PlayerData, 0, len(players))
	for _, player := range players {
		playerList = append(playerList, player)
	}
	return playerList
}

func RemovePlayer(playerID string) {
	playerMutex.Lock()
	defer playerMutex.Unlock()
	delete(players, playerID)
}

// アイテム関連
func SetItem(item domain.ItemData) {
	itemMutex.Lock()
	defer itemMutex.Unlock()
	items[item.ID] = item
}

func GetAllItems() []domain.ItemData {
	itemMutex.Lock()
	defer itemMutex.Unlock()
	itemList := make([]domain.ItemData, 0, len(items))
	for _, item := range items {
		itemList = append(itemList, item)
	}
	return itemList
}

func RemoveItem(itemID string) {
	itemMutex.Lock()
	defer itemMutex.Unlock()
	delete(items, itemID)
}

// 弾丸関連
func SetBullet(bullet domain.BulletData) {
	bulletMutex.Lock()
	defer bulletMutex.Unlock()
	bullets[bullet.ID] = bullet
}

func GetAllBullets() []domain.BulletData {
	bulletMutex.Lock()
	defer bulletMutex.Unlock()
	bulletList := make([]domain.BulletData, 0, len(bullets))
	for _, bullet := range bullets {
		bulletList = append(bulletList, bullet)
	}
	return bulletList
}

// クライアントIDとプレイヤーIDの紐付け関連
func SetClientPlayerID(clientID, playerID string) {
	if _, exists := GetPlayer(playerID); !exists {
		return
	}
	clientPlayerMutex.Lock()
	defer clientPlayerMutex.Unlock()
	clientIDToPlayerID[clientID] = playerID
}

func GetPlayerIDByClientID(clientID string) (string, bool) {
	//clientPlayerMutex.Lock()
	//defer clientPlayerMutex.Unlock()
	playerID, exists := clientIDToPlayerID[clientID]
	return playerID, exists
}

func GetClientIDByPlayerID(playerID string) (string, bool) {
	//clientPlayerMutex.Lock()
	//defer clientPlayerMutex.Unlock()
	for clientID, pID := range clientIDToPlayerID {
		if pID == playerID {
			return clientID, true
		}
	}
	return "", false
}

func RemoveClientPlayerID(clientID string) {
	clientPlayerMutex.Lock()
	defer clientPlayerMutex.Unlock()
	delete(clientIDToPlayerID, clientID)
}
