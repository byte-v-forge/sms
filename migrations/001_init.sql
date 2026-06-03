DROP TABLE IF EXISTS sms_route_offers;
DROP TABLE IF EXISTS sms_route_profiles;
DROP TABLE IF EXISTS sms_activations;
DROP TABLE IF EXISTS sms_order_codes;
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

DO $$
DECLARE
  primary_key_name text;
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_name = 'sms_provider_configs'
      AND column_name = 'provider_config_id'
  ) THEN
    UPDATE sms_provider_configs
    SET provider_key = provider_config_id
    WHERE (provider_key IS NULL OR provider_key = '')
      AND provider_config_id <> '';
  END IF;

  DELETE FROM sms_provider_configs
  WHERE ctid IN (
    SELECT ctid
    FROM (
      SELECT ctid,
        row_number() OVER (
          PARTITION BY provider_key
          ORDER BY updated_at DESC, created_at DESC, ctid DESC
        ) AS row_number
      FROM sms_provider_configs
    ) duplicate_rows
    WHERE row_number > 1
  );

  DELETE FROM sms_provider_configs
  WHERE provider_key IS NULL OR provider_key = '';

  SELECT conname
  INTO primary_key_name
  FROM pg_constraint
  WHERE conrelid = 'sms_provider_configs'::regclass
    AND contype = 'p';

  IF primary_key_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE sms_provider_configs DROP CONSTRAINT %I', primary_key_name);
  END IF;

  ALTER TABLE sms_provider_configs ALTER COLUMN provider_key SET NOT NULL;
  ALTER TABLE sms_provider_configs ADD CONSTRAINT sms_provider_configs_pkey PRIMARY KEY (provider_key);
  ALTER TABLE sms_provider_configs DROP COLUMN IF EXISTS provider_config_id;
  ALTER TABLE sms_provider_configs DROP COLUMN IF EXISTS display_name;
  ALTER TABLE sms_provider_configs DROP COLUMN IF EXISTS config;
END $$;

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

CREATE TABLE IF NOT EXISTS sms_order_codes (
  order_id text NOT NULL REFERENCES sms_orders(order_id) ON DELETE CASCADE,
  code_secret_id text NOT NULL,
  message_text text NOT NULL DEFAULT '',
  received_at timestamptz NOT NULL,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (order_id, code_secret_id, received_at)
);
CREATE INDEX IF NOT EXISTS idx_sms_order_codes_order_received
  ON sms_order_codes(order_id, received_at DESC);

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
