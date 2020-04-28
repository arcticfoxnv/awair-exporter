package main

import (
	"fmt"
	"github.com/arcticfoxnv/awair-exporter/awair"
	"github.com/arcticfoxnv/awair_api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Printf("awair-exporter v%s-%s", Version, Commit)
	log.Printf("-- awair_api v%s", awair_api.Version)

	config, err := loadConfig()
	if err != nil {
		log.Println("Failed to load config file:", err)
	}

	if err := preflightCheck(config); err != nil {
		log.Fatalln(err)
	}

	client := awair.NewClient(
		config.GetString(CFG_ACCESS_TOKEN),
		0,
		func(c *awair_api.Client) {
			c.UserAgent = fmt.Sprintf("awair-exporter/%s (https://github.com/arcticfoxnv/awair-exporter)", Version)
		},
	)

	userInfo, err := client.GetUserInfo()
	if err != nil {
		log.Fatalln("Failed to retrieve user info:", err)
	}

	config.SetDefault(CFG_TIER_NAME, userInfo.Tier)
	tierName := config.GetString(CFG_TIER_NAME)
	log.Println("API tier level:", tierName)

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
