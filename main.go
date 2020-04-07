package main

import (
	"awair-exporter/awair"
	"github.com/arcticfoxnv/awair_api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	log.Printf("awair-exporter v%s-%s", Version, Commit)
	log.Printf("-- awair_api v%s", awair_api.Version)
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

	client := awair.NewClient(
		config.AccessToken,
		0,
		func(c *awair_api.Client) {
			c.UserAgent = "awair-exporter (https://github.com/arcticfoxnv/awair-exporter)"
		},
	)

	userInfo, err := client.GetUserInfo()
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
	client.SetCacheTTL(cacheTTL)

	registry := prometheus.NewRegistry()
	registry.MustRegister(NewAwairCollector(client))

	e := NewExporterHTTP(client)
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	m.HandleFunc("/meta/usage", e.serveUsage)
	s := &http.Server{Addr: ":8080", Handler: m}

	log.Println("Starting HTTP listener on", s.Addr)
	s.ListenAndServe()
}
