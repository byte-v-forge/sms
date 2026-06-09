package grpcadapter

import (
	smsv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/contracts/sms/v1"
	"github.com/byte-v-forge/sms/internal/core"
)

func toProtoCatalogApplications(applications []core.CatalogApplication) []*smsv1.SmsApplicationInfo {
	out := make([]*smsv1.SmsApplicationInfo, 0, len(applications))
	for _, item := range applications {
		out = append(out, &smsv1.SmsApplicationInfo{
			ApplicationKey: item.ApplicationKey,
			DisplayName:    item.DisplayName,
		})
	}
	return out
}

func toProtoCatalogCountries(countries []core.CatalogCountry) []*smsv1.SmsCountry {
	out := make([]*smsv1.SmsCountry, 0, len(countries))
	for _, item := range countries {
		out = append(out, &smsv1.SmsCountry{
			CountryIso2:        item.CountryISO2,
			Name:               item.Name,
			CountryCallingCode: item.CountryCallingCode,
		})
	}
	return out
}
