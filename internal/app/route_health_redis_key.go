package app

import (
	"crypto/sha256"
	"encoding/hex"
)

func routeHealthRedisKey(kind string, routeKey string) string {
	digest := sha256.Sum256([]byte(routeKey))
	return "sms:route_health:" + kind + ":" + hex.EncodeToString(digest[:])
}
