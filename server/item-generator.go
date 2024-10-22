package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github/pikachu0310/hackathon24spring-server/domain"
)

const apiURL = "https://online-data.trap.games/api/items"

func init() {
	go generateItems()
}

func generateRandomPosition() domain.Vector3 {
	return domain.Vector3{
		X: rand.Float32()*20 - 10,
		Y: rand.Float32()*20 - 10,
		Z: 0,
	}
}

func fetchItemFromAPI() (domain.ItemData, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return domain.ItemData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.ItemData{}, fmt.Errorf("failed to fetch item: %s", resp.Status)
	}

	var item domain.ItemData
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return domain.ItemData{}, err
	}

	item.Type = "item"
	item.Position = generateRandomPosition()
	item.Rotation = 0
	size := (rand.Float32() + 0.5) * 0.25
	item.Size = size
	item.Mass = size * 5
	item.Speed = domain.Vector3{X: 0, Y: 0, Z: 0}

	return item, nil
}

func generateItems() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// すべてのアイテムの中で、lastTouchedが""なものだけ配信する
		for _, item := range GetAllItems() {
			if item.LastTouched == "" {
				broadcastItemUpdate(item)
			}
		}

		// アイテムが5つ以上あれば生成しない
		clientMutex.Lock()
		if len(items) >= 5 {
			clientMutex.Unlock()
			continue
		}
		clientMutex.Unlock()

		item, err := fetchItemFromAPI()
		if err != nil {
			log.Printf("Error fetching item: %v", err)
			continue
		}

		// 状態管理ファイルにアイテムをセット
		SetItem(item)

		// アイテム情報をブロードキャスト
		broadcastItemUpdate(item)
	}
}
