package domain

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type PlayerData struct {
	Type       string  `json:"type"`
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Position   Vector3 `json:"position"`
	Speed      Vector3 `json:"speed"`
	Rotation   float32 `json:"rotation"`
	MaxHP      float32 `json:"maxHP"`
	HP         float32 `json:"hp"`
	Mass       float32 `json:"mass"`
	Bounciness float32 `json:"bounciness"`
	Friction   float32 `json:"friction"`
	Size       float32 `json:"size"`
	Score      int     `json:"score"`
	GrabTarget string  `json:"grabTarget"`
	GrabLength float32 `json:"grabLength"`
	KillerID   string  `json:"killerID"`
}

type ItemData struct {
	Type        string  `json:"type"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Attribute   int     `json:"attribute"`
	Rarity      int     `json:"rarity"`
	Position    Vector3 `json:"position"`
	Speed       Vector3 `json:"speed"`
	Rotation    float32 `json:"rotation"`
	Size        float32 `json:"size"`
	Mass        float32 `json:"mass"`
}

type BulletData struct {
	Type          string  `json:"type"`
	ID            string  `json:"id"`
	FirstPosition Vector3 `json:"firstPosition"`
	FirstSpeed    Vector3 `json:"firstSpeed"`
	RemainingTime float32 `json:"remainingTime"`
	Attack        float32 `json:"attack"`
	Size          float32 `json:"size"`
	ShooterID     string  `json:"shooterID"`
}

type KillData struct {
	KilledPlayerID string `json:"killedPlayerID"`
	KillerPlayerID string `json:"killerPlayerID"`
}
