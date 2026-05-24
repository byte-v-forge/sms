package app

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"
)

type RandomIDGenerator struct{}

var fallbackIDCounter atomic.Uint64

func (RandomIDGenerator) NewID(prefix string) string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err == nil {
		return prefix + hex.EncodeToString(b[:])
	}
	binary.BigEndian.PutUint64(b[:8], uint64(time.Now().UnixNano()))
	binary.BigEndian.PutUint32(b[8:], uint32(fallbackIDCounter.Add(1)))
	return prefix + hex.EncodeToString(b[:])
}
