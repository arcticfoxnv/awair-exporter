package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCacheTTLByTierDefault(t *testing.T) {
	assert.Equal(t, CACHE_TIER_HOBBYIST, GetCacheTTLByTier("xxx"))
}

func TestGetCacheTTLByTierHobbyist(t *testing.T) {
	assert.Equal(t, CACHE_TIER_HOBBYIST, GetCacheTTLByTier("Hobbyist"))
}

func TestGetCacheTTLByTierEnterprise(t *testing.T) {
	assert.Equal(t, CACHE_TIER_ENTERPRISE, GetCacheTTLByTier("Enterprise"))
}
