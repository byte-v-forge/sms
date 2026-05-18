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
	"github.com/byte-v-forge/sms/internal/providers/herosms"
)

func main() {
	apiKey := flag.String("api-key", os.Getenv("HERO_SMS_API_KEY"), "HeroSMS API key. Can also use HERO_SMS_API_KEY.")
	endpoint := flag.String("endpoint", herosms.DefaultEndpoint, "HeroSMS handler API endpoint.")
	countryFilter := flag.String("country", "", "Optional HeroSMS country id filter, for example 6.")
	serviceFilter := flag.String("service", "", "Optional HeroSMS service code filter, for example tg.")
	minCount := flag.Int("min-count", 1, "Minimum available count to include.")
	limit := flag.Int("limit", 100, "Maximum rows to print. Use 0 for all rows.")
	format := flag.String("format", "table", "Output format: table or json.")
	flag.Parse()

	if strings.TrimSpace(*apiKey) == "" {
		fmt.Fprintln(os.Stderr, "missing HeroSMS API key: pass -api-key or set HERO_SMS_API_KEY")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := catalog.List(ctx, catalog.ProviderConfig{
		ProviderKey: herosms.ProviderKey,
		APIKey:      *apiKey,
		Endpoint:    *endpoint,
	}, catalog.ListOptions{
		ServiceKey: *serviceFilter,
		CountryID:  *countryFilter,
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

	fmt.Printf("HeroSMS balance: %s\n", result.Balance.AmountDecimal)
	fmt.Printf("Rows: %d", len(result.Options))
	if *countryFilter != "" || *serviceFilter != "" {
		fmt.Printf(" (country=%q service=%q)", *countryFilter, *serviceFilter)
	}
	fmt.Println()
	fmt.Println("country_id\tcountry\tservice_code\tservice\tcost\tcount\tphysical\troute")
	for _, option := range result.Options {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%d\t%d\tprovider=%s application=%s upstream_service=%s provider_country=%s\n",
			option.ProviderCountryID,
			firstNonEmpty(option.CountryNameCN, option.CountryName),
			option.UpstreamServiceKey,
			option.ServiceName,
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
