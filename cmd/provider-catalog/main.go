package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/byte-v-forge/sms/internal/catalog"
)

func main() {
	provider := flag.String("provider", "", "Provider key: 5sim, herosms, or smsbower.")
	apiKey := flag.String("api-key", "", "Provider API key/token. Can also use provider-specific env vars.")
	endpoint := flag.String("endpoint", "", "Optional provider API endpoint override.")
	currency := flag.String("currency", "", "Optional currency code, mainly for 5sim.")
	country := flag.String("country", "", "Optional provider country id filter.")
	service := flag.String("service", "", "Optional provider service/product code filter.")
	operator := flag.String("operator", "", "Optional 5sim operator filter.")
	minCount := flag.Int("min-count", 1, "Minimum available count to include.")
	limit := flag.Int("limit", 50, "Maximum options to print. Use 0 for all options.")
	format := flag.String("format", "table", "Output format: table or json.")
	flag.Parse()

	key := firstNonEmpty(*apiKey, providerAPIKeyEnv(*provider))
	if strings.TrimSpace(*provider) == "" {
		fmt.Fprintln(os.Stderr, "missing provider: pass -provider 5sim|herosms|smsbower")
		os.Exit(2)
	}
	if key == "" {
		fmt.Fprintln(os.Stderr, "missing provider API key: pass -api-key or set provider env var")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	result, err := catalog.List(ctx, catalog.ProviderConfig{
		ProviderKey:  *provider,
		APIKey:       key,
		Endpoint:     *endpoint,
		CurrencyCode: *currency,
	}, catalog.ListOptions{
		ServiceKey: *service,
		CountryID:  *country,
		Operator:   *operator,
		MinCount:   *minCount,
		Limit:      *limit,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *format == "json" {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("Provider: %s\n", result.ProviderKey)
	fmt.Printf("Balance: %s %s\n", result.Balance.AmountDecimal, result.Balance.CurrencyCode)
	fmt.Printf("Countries: %d Services: %d Options: %d\n", len(result.Countries), len(result.Services), len(result.Options))
	fmt.Println("country_id\tcountry\tservice_code\tservice\tprovider_or_operator\tcost\tcount\tphysical\troute")
	for _, option := range result.Options {
		upstream := firstNonEmpty(option.UpstreamProviderID, option.UpstreamOperator)
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\tprovider=%s application=%s upstream_service=%s provider_country=%s\n",
			option.ProviderCountryID,
			firstNonEmpty(option.CountryNameCN, option.CountryName),
			option.UpstreamServiceKey,
			option.ServiceName,
			upstream,
			option.Price.AmountDecimal,
			option.AvailableCount,
			option.PhysicalCount,
			option.ProviderKey,
			option.ApplicationKey,
			option.UpstreamServiceKey,
			option.ProviderCountryID,
		)
	}
}

func providerAPIKeyEnv(provider string) string {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "5sim":
		return os.Getenv("FIVESIM_TOKEN")
	case "herosms":
		return os.Getenv("HERO_SMS_API_KEY")
	case "smsbower":
		return os.Getenv("SMSBOWER_API_KEY")
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
