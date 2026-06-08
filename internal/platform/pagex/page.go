package pagex

const (
	DefaultLimit = 100
	MaxLimit     = 500
)

func NormalizeLimit(limit int, defaultLimit int, maxLimit int) int {
	if defaultLimit <= 0 {
		defaultLimit = DefaultLimit
	}
	if maxLimit <= 0 {
		maxLimit = defaultLimit
	}
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}
