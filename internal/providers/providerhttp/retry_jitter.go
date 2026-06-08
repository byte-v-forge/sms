package providerhttp

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

func randomJitter(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	var seed [8]byte
	if _, err := rand.Read(seed[:]); err != nil {
		return max / 2
	}
	return time.Duration(binary.LittleEndian.Uint64(seed[:]) % uint64(max+1))
}
