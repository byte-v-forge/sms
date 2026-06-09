package app

func priceLess(left routeCandidate, right routeCandidate) (bool, bool) {
	if left.hasPrice != right.hasPrice {
		return left.hasPrice, true
	}
	if left.hasPrice {
		if cmp := left.price.Cmp(right.price); cmp != 0 {
			return cmp < 0, true
		}
	}
	return false, false
}
