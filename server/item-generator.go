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

const apiURL = "https://online-data.trap.games/"

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
	//item.Rotation = rand.Float32() * 360
	item.Rotation = 0
	size := (rand.Float32() + 0.5) * 0.25
	item.Size = size
	item.Mass = size * 5
	item.Speed = domain.Vector3{X: 0, Y: 0, Z: 0} // 初期速度はゼロ

	return item, nil
}

func generateItems() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mutex.Lock()
		if len(items) >= 5 {
			mutex.Unlock()
			continue
		}
		mutex.Unlock()

		item, err := fetchItemFromAPI()
		if err != nil {
			log.Printf("Error fetching item: %v", err)
			continue
		}

		mutex.Lock()
		items[item.ID] = item
		mutex.Unlock()

		broadcastItems()
	}
}

func broadcastItems() {
	mutex.Lock()
	defer mutex.Unlock()

	itemList := make([]domain.ItemData, 0, len(items))
	for _, item := range items {
		itemList = append(itemList, item)
	}
	itemDataJSON, err := json.Marshal(itemList)
	if err != nil {
		log.Printf("Error marshalling item data: %v", err)
		return
	}

	for client := range clients {
		client.send <- itemDataJSON
	}
}
