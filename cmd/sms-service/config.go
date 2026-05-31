package main

import (
	"log"

	"github.com/byte-v-forge/common-lib/envx"
	"github.com/byte-v-forge/common-lib/natseventbus"
)

type config struct {
	ListenAddr         string
	PGDSN              string
	HTTPTimeoutSeconds int
	ProviderHTTPProxy  string
	DashboardHTTPAddr  string
	DashboardStaticDir string
	PlatformNATSURL    string
	PlatformRedisURL   string
	EventStreamName    string
}

func loadConfig() config {
	cfg := config{
		ListenAddr:         envx.StringDefault("SMS_LISTEN_ADDR", ":50051"),
		PGDSN:              envx.StringDefault("SMS_PG_DSN", envx.StringDefault("PG_DSN", "")),
		HTTPTimeoutSeconds: mustEnvInt("SMS_HTTP_TIMEOUT_SECONDS", 20),
		ProviderHTTPProxy:  envx.StringDefault("SMS_PROVIDER_HTTP_PROXY", ""),
		DashboardHTTPAddr:  envx.StringDefault("SMS_DASHBOARD_HTTP_ADDR", ":8080"),
		DashboardStaticDir: envx.StringDefault("SMS_DASHBOARD_STATIC_DIR", "/app/dashboard/sms"),
		PlatformNATSURL:    envx.StringDefault("PLATFORM_NATS_URL", ""),
		PlatformRedisURL:   envx.StringDefault("PLATFORM_REDIS_URL", ""),
		EventStreamName:    envx.StringDefault("PLATFORM_EVENT_STREAM_NAME", natseventbus.DefaultStream),
	}
	if cfg.HTTPTimeoutSeconds <= 0 {
		cfg.HTTPTimeoutSeconds = 20
	}
	return cfg
}

func mustEnvInt(name string, fallback int) int {
	parsed, err := envx.IntStrict(name, fallback)
	if err != nil {
		log.Fatal(err)
	}
	return parsed
}
