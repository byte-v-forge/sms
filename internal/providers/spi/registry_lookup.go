package spi

func (r *Registry) Get(providerKey string) (Plugin, bool) {
	if r == nil {
		return nil, false
	}
	plugin, ok := r.plugins[NormalizeKey(providerKey)]
	return plugin, ok
}

func (r *Registry) Supports(providerKey string) bool {
	_, ok := r.Get(providerKey)
	return ok
}
