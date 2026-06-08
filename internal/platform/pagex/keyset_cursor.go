package pagex

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type KeysetCursor struct {
	UpdatedAt time.Time `json:"updated_at"`
	ID        string    `json:"id"`
}

func EncodeKeysetCursor(updatedAt time.Time, id string) string {
	id = strings.TrimSpace(id)
	if updatedAt.IsZero() || id == "" {
		return ""
	}
	payload := KeysetCursor{UpdatedAt: updatedAt.UTC(), ID: id}
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

func DecodeKeysetCursor(value string) (KeysetCursor, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return KeysetCursor{}, nil
	}
	data, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return KeysetCursor{}, fmt.Errorf("invalid page cursor")
	}
	var cursor KeysetCursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return KeysetCursor{}, fmt.Errorf("invalid page cursor")
	}
	if cursor.UpdatedAt.IsZero() || strings.TrimSpace(cursor.ID) == "" {
		return KeysetCursor{}, fmt.Errorf("invalid page cursor")
	}
	cursor.UpdatedAt = cursor.UpdatedAt.UTC()
	cursor.ID = strings.TrimSpace(cursor.ID)
	return cursor, nil
}
