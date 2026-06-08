package pagex

import (
	"strings"
	"time"
)

type KeysetPage[T any] struct {
	Items      []T
	NextCursor string
	HasMore    bool
}

type KeysetCursorOf[T any] func(T) KeysetCursor

func KeysetLookaheadLimit(limit int) int {
	return NormalizePageLimit(limit) + 1
}

func HasKeysetCursor(cursor KeysetCursor) bool {
	return !cursor.UpdatedAt.IsZero() && strings.TrimSpace(cursor.ID) != ""
}

func NewKeysetPage[T any](rows []T, limit int, cursorOf KeysetCursorOf[T]) KeysetPage[T] {
	limit = NormalizePageLimit(limit)
	items, hasMore := TrimLimit(rows, limit)
	page := KeysetPage[T]{Items: items, HasMore: hasMore}
	if hasMore && len(items) > 0 && cursorOf != nil {
		page.NextCursor = EncodeKeysetCursorValue(cursorOf(items[len(items)-1]))
	}
	return page
}

func EncodeKeysetCursorValue(cursor KeysetCursor) string {
	return EncodeKeysetCursor(cursor.UpdatedAt, cursor.ID)
}

func KeysetCursorValue(updatedAt time.Time, id string) KeysetCursor {
	id = strings.TrimSpace(id)
	if updatedAt.IsZero() || id == "" {
		return KeysetCursor{}
	}
	return KeysetCursor{UpdatedAt: updatedAt.UTC(), ID: id}
}
