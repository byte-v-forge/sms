package core

import "time"

type RouteOffer struct {
	ProviderKey             string
	ProviderDisplayName     string
	UpstreamProviderID      string
	UpstreamProviderName    string
	ApplicationKey          string
	ApplicationName         string
	CountryISO2             string
	CountryName             string
	CountryCallingCode      string
	Price                   Money
	AvailableCount          int
	SupportsCancel          bool
	SupportsAdditionalCode  bool
	RequiresMarkMessageSent bool
	ObservedAt              time.Time
	Route                   Route
}
