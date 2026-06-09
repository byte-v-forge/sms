package fivesim

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/byte-v-forge/sms/internal/core"
)

func acquireNumberPath(route core.Route) (string, error) {
	country := strings.TrimSpace(route.ProviderCountryID)
	if country == "" {
		return "", core.NewError(core.CodeValidationFailed, "5sim country is required", false)
	}
	product := strings.TrimSpace(route.UpstreamServiceKey)
	if product == "" {
		return "", core.NewError(core.CodeValidationFailed, "5sim product is required", false)
	}
	operator, err := operatorForRoute(route)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/v1/user/buy/activation/%s/%s/%s", url.PathEscape(country), url.PathEscape(operator), url.PathEscape(product)), nil
}
