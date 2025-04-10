package event

import (
	"ptm/internal/db/redis"
	"ptm/internal/models"
	"time"
)

func AppendEvent(stream string, event models.Event) error {
	data := map[string]interface{}{
		"entity_id": event.EntityID,
		"type":      event.Type,
		"payload":   event.Payload,
		"timestamp": event.Timestamp.Format(time.RFC3339),
	}
	return redis.AppendEventToStream(stream, data)
}

func RebuildState(stream, entityID string) ([]models.Event, error) {
	var events []models.Event
	res, err := redis.ReadEventsFromStream(stream)
	if err != nil {
		return nil, err
	}

	for _, msg := range res {
		if msg.Values["entity_id"] == entityID {
			payload, _ := msg.Values["payload"].(string)
			events = append(events, models.Event{
				ID:        msg.ID,
				EntityID:  msg.Values["entity_id"].(string),
				Type:      msg.Values["type"].(string),
				Payload:   payload,
				Timestamp: time.Now(),
			})
		}
	}
	return events, nil
}
