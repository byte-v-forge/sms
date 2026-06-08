package main

import (
	"github.com/byte-v-forge/sms/internal/platform/envx"
	"github.com/byte-v-forge/sms/internal/platform/natseventbus"
)

type config struct {
	ListenAddr         string
	PGDSN              string
	HTTPTimeoutSeconds int
	ProviderHTTPProxy  string
	DashboardHTTPAddr  string
	DashboardStaticDir string
	NATSURL            string
	RedisURL           string
	EventStreamName    string
}

func loadConfig() (config, error) {
	httpTimeoutSeconds, err := envx.IntStrict("SMS_HTTP_TIMEOUT_SECONDS", 20)
	if err != nil {
		return config{}, err
	}
	cfg := config{
		ListenAddr:         envx.StringDefault("SMS_LISTEN_ADDR", ":50051"),
		PGDSN:              envx.StringDefault("SMS_PG_DSN", ""),
		HTTPTimeoutSeconds: httpTimeoutSeconds,
		ProviderHTTPProxy:  envx.StringDefault("SMS_PROVIDER_HTTP_PROXY", ""),
		DashboardHTTPAddr:  envx.StringDefault("SMS_DASHBOARD_HTTP_ADDR", ":8080"),
		DashboardStaticDir: envx.StringDefault("SMS_DASHBOARD_STATIC_DIR", "/app/dashboard/sms"),
		NATSURL:            envx.StringDefault("SMS_NATS_URL", ""),
		RedisURL:           envx.StringDefault("SMS_REDIS_URL", ""),
		EventStreamName:    envx.StringDefault("SMS_EVENT_STREAM_NAME", natseventbus.DefaultStream),
	}
	if cfg.HTTPTimeoutSeconds <= 0 {
		cfg.HTTPTimeoutSeconds = 20
	}
	return cfg, nil
}
