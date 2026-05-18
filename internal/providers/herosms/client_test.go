package herosms

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/byte-v-forge/sms/internal/core"
)

func TestAcquireNumberParsesAccessNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("action"); got != "getNumber" {
			t.Fatalf("action = %q, want getNumber", got)
		}
		if got := r.URL.Query().Get("service"); got != "go" {
			t.Fatalf("service = %q, want go", got)
		}
		_, _ = w.Write([]byte("ACCESS_NUMBER:123:628123456789"))
	}))
	defer server.Close()

	client, err := New(Config{Endpoint: server.URL, APIKey: "test-token"}, server.Client())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	resp, err := client.AcquireNumber(context.Background(), core.ProviderAcquireRequest{
		Route:  core.Route{UpstreamServiceKey: "go", ProviderCountryID: "6"},
		Target: core.Target{ApplicationKey: "gojek", CountryISO2: "ID", CountryCallingCode: "62"},
	})
	if err != nil {
		t.Fatalf("AcquireNumber() error = %v", err)
	}
	if resp.UpstreamActivationID != "123" {
		t.Fatalf("activation id = %q", resp.UpstreamActivationID)
	}
	if resp.PhoneNumber.E164 != "+628123456789" || resp.PhoneNumber.NationalNumber != "8123456789" {
		t.Fatalf("phone = %#v", resp.PhoneNumber)
	}
}

func TestSetStatusCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("status"); got != "8" {
			t.Fatalf("status = %q, want 8", got)
		}
		_, _ = w.Write([]byte("ACCESS_CANCEL"))
	}))
	defer server.Close()

	client, err := New(Config{Endpoint: server.URL, APIKey: "test-token"}, server.Client())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if err := client.SetStatus(context.Background(), "123", core.ActionCancelActivation); err != nil {
		t.Fatalf("SetStatus() error = %v", err)
	}
}

func TestPolicyRequiresTwoMinutesBeforeCancel(t *testing.T) {
	client, err := New(Config{Endpoint: "https://example.test/api", APIKey: "test-token"}, nil)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got := client.Policy().CancelAllowedAfter; got != 2*time.Minute {
		t.Fatalf("CancelAllowedAfter = %s, want 2m", got)
	}
}

func TestListCatalogData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("action") {
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

	client, err := New(Config{Endpoint: server.URL, APIKey: "test-token"}, server.Client())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	countries, err := client.ListCountries(context.Background())
	if err != nil {
		t.Fatalf("ListCountries() error = %v", err)
	}
	if len(countries) != 1 || countries[0].CountryID != "6" || countries[0].NameCN != "印度尼西亚" || !countries[0].Retry {
		t.Fatalf("countries = %#v", countries)
	}
	services, err := client.ListApplications(context.Background())
	if err != nil {
		t.Fatalf("ListApplications() error = %v", err)
	}
	if len(services) != 1 || services[0].UpstreamServiceKey != "tg" || services[0].DisplayName != "Telegram" {
		t.Fatalf("services = %#v", services)
	}
	offers, err := client.ListPriceOffers(context.Background(), "tg", "6")
	if err != nil {
		t.Fatalf("ListPriceOffers() error = %v", err)
	}
	if len(offers) != 1 || offers[0].CountryID != "6" || offers[0].UpstreamServiceKey != "tg" || offers[0].Price.AmountDecimal != "0.15" || offers[0].AvailableCount != 24794 || offers[0].PhysicalCount != 22846 {
		t.Fatalf("offers = %#v", offers)
	}
}
