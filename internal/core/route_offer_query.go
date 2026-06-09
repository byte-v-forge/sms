package core

type RouteOfferQuery struct {
	ApplicationKey     string
	CountryISO2        string
	CountryCallingCode string
	ProviderKeys       []string
}
