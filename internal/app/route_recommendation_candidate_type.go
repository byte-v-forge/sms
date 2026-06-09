package app

import (
	"math/big"

	"github.com/byte-v-forge/sms/internal/core"
)

type routeCandidate struct {
	offer    core.RouteOffer
	price    *big.Rat
	hasPrice bool
	score    int32
}
