package core

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
