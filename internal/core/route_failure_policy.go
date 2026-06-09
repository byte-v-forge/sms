package core

import "time"

type RouteFailurePolicy struct {
	ScopeKey         string
	FailureThreshold int
	FailureWindow    time.Duration
	DisableTTL       time.Duration
}
