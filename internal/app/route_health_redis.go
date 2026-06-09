package app

func (s *RedisRouteHealthStore) available() bool {
	return s != nil && s.client != nil
}
