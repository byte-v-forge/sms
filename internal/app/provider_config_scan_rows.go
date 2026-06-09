package app

import (
	smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"
	"github.com/jackc/pgx/v5"
)

func scanProviderConfigs(rows pgx.Rows) ([]*smsinternalv1.SmsProviderConfig, error) {
	defer rows.Close()
	configs := []*smsinternalv1.SmsProviderConfig{}
	for rows.Next() {
		config, err := scanProviderConfig(rows)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}
	return configs, rows.Err()
}
