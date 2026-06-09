package herosms

func activationOfferPurchaseTiers(offer activationOffer) []activationOfferPurchaseTier {
	if tiers := activationOfferPriceMapTiers(offer.PriceMap); len(tiers) > 0 {
		return tiers
	}
	return activationOfferDefaultTier(offer)
}
