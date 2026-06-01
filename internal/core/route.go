package core

import "time"

type Route struct {
	ProviderKey        string
	ApplicationKey     string
	UpstreamServiceKey string
	CountryISO2        string
	CountryCallingCode string
	MinAvailableCount  int
	MinPrice           Money
	MaxPrice           Money
	ProviderCountryID  string
	UpstreamProviderID string
	FailurePolicy      RouteFailurePolicy
}

type RouteFailurePolicy struct {
	ScopeKey         string
	FailureThreshold int
	FailureWindow    time.Duration
	DisableTTL       time.Duration
}

type RouteOfferQuery struct {
	ApplicationKey     string
	CountryISO2        string
	CountryCallingCode string
	ProviderKey        string
}

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
