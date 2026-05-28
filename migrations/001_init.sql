DROP TABLE IF EXISTS sms_route_offers;
DROP TABLE IF EXISTS sms_route_profiles;
DROP TABLE IF EXISTS sms_activations;
DROP TABLE IF EXISTS sms_orders;

CREATE TABLE IF NOT EXISTS sms_provider_configs (
  provider_key text PRIMARY KEY,
  enabled boolean NOT NULL DEFAULT true,
  credential_secret text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_sms_provider_configs_enabled
  ON sms_provider_configs(enabled, provider_key);

CREATE TABLE IF NOT EXISTS sms_orders (
  order_id text PRIMARY KEY,
  request_id text NOT NULL DEFAULT '',
  provider_key text NOT NULL DEFAULT '',
  upstream_order_id text NOT NULL DEFAULT '',
  target_application_key text NOT NULL DEFAULT '',
  target_country_iso2 text NOT NULL DEFAULT '',
  target_country_calling_code text NOT NULL DEFAULT '',
  phone_e164 text NOT NULL DEFAULT '',
  status text NOT NULL DEFAULT '',
  price_currency text NOT NULL DEFAULT '',
  price_amount text NOT NULL DEFAULT '',
  acquired_at timestamptz,
  expires_at timestamptz,
  updated_at timestamptz NOT NULL DEFAULT now(),
  cancel_allowed_at timestamptz,
  last_error_code text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sms_orders_request_id
  ON sms_orders(request_id) WHERE request_id <> '';
CREATE INDEX IF NOT EXISTS idx_sms_orders_status_updated
  ON sms_orders(status, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_sms_orders_provider_upstream
  ON sms_orders(provider_key, upstream_order_id);

CREATE TABLE IF NOT EXISTS sms_platform_event_outbox (
  event_id TEXT PRIMARY KEY,
  subject TEXT NOT NULL,
  event_name TEXT NOT NULL,
  idempotency_key TEXT NOT NULL DEFAULT '',
  envelope BYTEA NOT NULL,
  status TEXT NOT NULL DEFAULT 'PENDING',
  attempt_count INT NOT NULL DEFAULT 0,
  next_attempt_at BIGINT NOT NULL DEFAULT 0,
  last_error TEXT NOT NULL DEFAULT '',
  published_at BIGINT NOT NULL DEFAULT 0,
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sms_platform_event_outbox_pending
  ON sms_platform_event_outbox(status, next_attempt_at, created_at);
