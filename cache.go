package main

import "time"

const (
	CACHE_TIER_HOBBYIST   = 5 * time.Minute
  CACHE_TIER_ENTERPRISE = time.Second
)

func GetCacheTTLByTier(tier string) time.Duration {
	switch tier {
	case "Hobbyist":
		return CACHE_TIER_HOBBYIST
  case "Enterprise":
    return CACHE_TIER_ENTERPRISE
	default:
		return CACHE_TIER_HOBBYIST
	}
}
