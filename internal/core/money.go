package core

type Money struct {
	CurrencyCode  string
	AmountDecimal string
}

type PhoneNumber struct {
	E164               string
	NationalNumber     string
	CountryISO2        string
	CountryCallingCode string
}

type Target struct {
	ApplicationKey     string
	CountryISO2        string
	CountryCallingCode string
}
