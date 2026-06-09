package core

type CatalogApplication struct {
	ApplicationKey string
	DisplayName    string
	Aliases        []string
}

type CatalogCountry struct {
	CountryISO2        string
	Name               string
	CountryCallingCode string
}

type CatalogApplicationQuery struct {
	ProviderKeys []string
	SearchText   string
}

type CatalogCountryQuery struct {
	ProviderKeys   []string
	ApplicationKey string
}
