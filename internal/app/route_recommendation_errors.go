package app

import (
	"fmt"

	"github.com/byte-v-forge/sms/internal/core"
)

func routeRecommendationUnavailableError(target core.Target) error {
	return core.NewError(
		core.CodeSupplyUnavailable,
		fmt.Sprintf("sms route currently unavailable for %s/%s/%s", target.ApplicationKey, target.CountryISO2, target.CountryCallingCode),
		true,
	)
}
