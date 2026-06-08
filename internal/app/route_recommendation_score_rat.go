package app

import (
	"math"
	"math/big"
)

func ratScore(value *big.Rat) int32 {
	if value == nil {
		return 0
	}
	score, _ := value.Float64()
	score = math.Round(score)
	if score < 0 {
		return 0
	}
	if score > 1000 {
		return 1000
	}
	return int32(score)
}
