package app

type CompositeOrderStore struct {
	active  OrderStore
	history OrderStore
}

func NewCompositeOrderStore(active OrderStore, history OrderStore) *CompositeOrderStore {
	return &CompositeOrderStore{active: active, history: history}
}
