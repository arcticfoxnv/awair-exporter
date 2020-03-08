package main

import (
	"awair-exporter/awair"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CACHE_TIER_HOBBYIST = 5 * time.Minute
)

func getEnvWithDefault(key, defaultValue string) string {
	if value, present := os.LookupEnv(key); present {
		return value
	}
	return defaultValue
}

func GetCacheTTLByTier(tier string) time.Duration {
	switch tier {
	case "Hobbyist":
		return CACHE_TIER_HOBBYIST
	default:
		return CACHE_TIER_HOBBYIST
	}
}

func main() {
	cfgFilename := getEnvWithDefault("AWAIR_CONFIG_FILE", "awair.toml")
	config, err := LoadConfig(cfgFilename)
	if err != nil {
		log.Println("Failed to load config file:", err)
		config = &Config{}
	}
	config.AccessToken = getEnvWithDefault("AWAIR_ACCESS_TOKEN", config.AccessToken)
	if config.AccessToken == "" {
		log.Fatalln("Cannot start exporter, access token is missing")
	}

	client := awair.NewClient(config.AccessToken)
	userInfo, err := client.UserInfo()
	if err != nil {
		log.Fatalln("Failed to retrieve user info:", err)
	}

	tierName := userInfo.Tier
	log.Println("API tier level:", tierName)
	if config.Tier != "" {
		tierName = config.Tier
	}

	cacheTTL := GetCacheTTLByTier(tierName)
	log.Printf("Setting cache key ttl to %d seconds\n", cacheTTL/time.Second)

	e := NewExporterHTTP(client, cacheTTL)
	m := http.NewServeMux()
	m.HandleFunc("/metrics", e.serveLatest)
	m.HandleFunc("/data/latest", e.serveLatest)
	m.HandleFunc("/meta/usage", e.serveUsage)
	s := &http.Server{Addr: ":8080", Handler: m}

	log.Println("Starting HTTP listener on", s.Addr)
	s.ListenAndServe()
}
