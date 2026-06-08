package spi

import (
	"fmt"
	"sort"
)

func NewRegistry(plugins ...Plugin) (*Registry, error) {
	registry := &Registry{plugins: make(map[string]Plugin, len(plugins))}
	for _, plugin := range plugins {
		if err := registry.add(plugin); err != nil {
			return nil, err
		}
	}
	sort.Strings(registry.keys)
	return registry, nil
}

func EmptyRegistry() *Registry {
	return &Registry{plugins: map[string]Plugin{}, keys: []string{}}
}

func (r *Registry) add(plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("sms provider plugin is required")
	}
	key := NormalizeKey(plugin.Key())
	if key == "" {
		return fmt.Errorf("sms provider plugin key is required")
	}
	if _, exists := r.plugins[key]; exists {
		return fmt.Errorf("duplicate sms provider plugin %q", key)
	}
	r.plugins[key] = plugin
	r.keys = append(r.keys, key)
	return nil
}
