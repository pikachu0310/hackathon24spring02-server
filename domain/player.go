package domain

type PlayerData struct {
	ID       string  `json:"id"`
	Position Vector3 `json:"position"`
	Speed    float32 `json:"speed"`
}

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}
