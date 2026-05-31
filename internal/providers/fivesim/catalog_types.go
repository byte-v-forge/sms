package fivesim

import "github.com/byte-v-forge/sms/internal/core"

type Country struct {
	CountryID          string
	Name               string
	CountryISO2        string
	CountryCallingCode string
}

type PriceOffer struct {
	CountryID          string
	UpstreamServiceKey string
	Operator           string
	Price              core.Money
	AvailableCount     int
	SuccessRate        float64
}
