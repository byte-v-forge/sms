package spi

import smsinternalv1 "github.com/byte-v-forge/sms/gen/go/byte/v/forge/sms/private/v1"

func (r *Registry) Descriptors() []*smsinternalv1.SmsProviderPluginDescriptor {
	if r == nil {
		return nil
	}
	out := make([]*smsinternalv1.SmsProviderPluginDescriptor, 0, len(r.keys))
	for _, key := range r.keys {
		plugin := r.plugins[key]
		out = append(out, &smsinternalv1.SmsProviderPluginDescriptor{
			ProviderKey:  plugin.Key(),
			DisplayName:  plugin.DisplayName(),
			Capabilities: plugin.Capabilities(),
		})
	}
	return out
}
