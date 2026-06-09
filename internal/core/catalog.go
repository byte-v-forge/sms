package core

type CatalogApplication struct {
	ApplicationKey string
	DisplayName    string
}

type CatalogCountry struct {
	CountryISO2        string
	Name               string
	CountryCallingCode string
}

type CatalogApplicationQuery struct {
	ProviderKeys []string
}

type CatalogCountryQuery struct {
	ProviderKeys   []string
	ApplicationKey string
}
