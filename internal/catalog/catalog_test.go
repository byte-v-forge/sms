package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/byte-v-forge/sms/internal/providers/herosms"
)

func TestListHeroSMSCatalogBuildsSelectableOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("action") {
		case "getBalance":
			_, _ = w.Write([]byte("ACCESS_BALANCE:1.363"))
		case "getCountries":
			_, _ = w.Write([]byte(`[{"id":6,"rus":"Индонезия","eng":"Indonesia","chn":"印度尼西亚","visible":1,"retry":1}]`))
		case "getServicesList":
			_, _ = w.Write([]byte(`{"status":"success","services":[{"code":"tg","name":"Telegram"}]}`))
		case "getPrices":
			if got := r.URL.Query().Get("country"); got != "6" {
				t.Fatalf("country = %q, want 6", got)
			}
			if got := r.URL.Query().Get("service"); got != "tg" {
				t.Fatalf("service = %q, want tg", got)
			}
			_, _ = w.Write([]byte(`{"6":{"tg":{"cost":0.15,"count":24794,"physicalCount":22846}}}`))
		default:
			t.Fatalf("unexpected action %q", r.URL.Query().Get("action"))
		}
	}))
	defer server.Close()

	result, err := List(context.Background(), ProviderConfig{
		ProviderKey: herosms.ProviderKey,
		APIKey:      "test-token",
		Endpoint:    server.URL,
	}, ListOptions{CountryID: "6", ServiceKey: "tg", Limit: 10})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if result.ProviderKey != herosms.ProviderKey || result.Balance.AmountDecimal != "1.363" {
		t.Fatalf("catalog header = %#v", result)
	}
	if len(result.Countries) != 1 || result.Countries[0].ProviderCountryID != "6" {
		t.Fatalf("countries = %#v", result.Countries)
	}
	if len(result.Services) != 1 || result.Services[0].DisplayName != "Telegram" {
		t.Fatalf("services = %#v", result.Services)
	}
	if len(result.Options) != 1 {
		t.Fatalf("options = %#v", result.Options)
	}
	option := result.Options[0]
	if option.ProviderKey != "herosms" || option.ApplicationKey != "tg" || option.UpstreamServiceKey != "tg" || option.ProviderCountryID != "6" {
		t.Fatalf("route option = %#v", option)
	}
	if option.CountryNameCN != "印度尼西亚" || option.ServiceName != "Telegram" || option.Price.AmountDecimal != "0.15" || option.AvailableCount != 24794 || option.PhysicalCount != 22846 {
		t.Fatalf("display option = %#v", option)
	}
}
