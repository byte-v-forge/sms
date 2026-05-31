package smsbower

import "github.com/byte-v-forge/sms/internal/core"

type ApplicationOffer struct {
	ApplicationKey     string
	UpstreamServiceKey string
	DisplayName        string
}

type Country struct {
	CountryID string
	Name      string
}

type PriceOffer struct {
	CountryID          string
	UpstreamServiceKey string
	ProviderID         string
	Price              core.Money
	AvailableCount     int
}
