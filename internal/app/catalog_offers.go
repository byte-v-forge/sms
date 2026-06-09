package app

import "github.com/byte-v-forge/sms/internal/core"

const catalogProviderConcurrency = 4

type RouteOfferList struct {
	Offers         []core.RouteOffer
	ProviderErrors []ProviderLookupError
}

type ProviderLookupError struct {
	ProviderKey         string
	ProviderDisplayName string
	Err                 error
}

type catalogProviderOffersResult struct {
	providerKey         string
	providerDisplayName string
	offers              []core.RouteOffer
	err                 error
}

type catalogProviderApplicationsResult struct {
	applications []core.CatalogApplication
	err          error
}

type catalogProviderCountriesResult struct {
	countries []core.CatalogCountry
	err       error
}
