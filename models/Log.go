package models

type Log struct {
	EntityType string `json:"entity_type"`
	EntityId   string `json:"entity_id"`
	Action     string `json:"action"`
	Details    string `json:"details"`
	CreatedAt  int64  `gorm:"column:created_at"`
}
