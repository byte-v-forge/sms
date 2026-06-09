package herosms

import "github.com/byte-v-forge/sms/internal/core"

type PriceOffer struct {
	CountryID          string
	UpstreamServiceKey string
	Operator           string
	OperatorName       string
	Price              core.Money
	AvailableCount     int
}

type countryMetadata struct {
	ID          string
	Name        string
	ISO2        string
	CallingCode string
}
