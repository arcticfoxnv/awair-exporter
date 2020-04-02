package main

import (
	awair "awair-exporter/awair/client"
	"awair-exporter/awair/structs"
	"fmt"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"time"
)

type ExporterHTTP struct {
	client      *awair.Client
	clientCache *cache.Cache
}

func NewExporterHTTP(client *awair.Client, cacheTTL time.Duration) *ExporterHTTP {
	return &ExporterHTTP{
		client:      client,
		clientCache: cache.New(cacheTTL, 10*time.Minute),
	}
}

// TODO: This should be merged with the collector version
// TODO: Ideally, the cache should also be part of a high-level client
func (e *ExporterHTTP) getDeviceList() (*structs.DeviceList, error) {
	if data, found := e.clientCache.Get(DEVICES_KEY); found {
		return data.(*structs.DeviceList), nil
	}
	log.Printf("Fetching device list")

	devices, err := e.client.Devices()
	if err != nil {
		return nil, err
	}

	e.clientCache.Set(DEVICES_KEY, devices, cache.DefaultExpiration)
	return devices, nil
}

func (e *ExporterHTTP) serveUsage(w http.ResponseWriter, r *http.Request) {
	userInfo, err := e.client.UserInfo()
	if err != nil {
		fmt.Printf("Error while getting user info: %s\n", err)
		return
	}
	fmt.Fprintf(w, "%+v\n", userInfo)

	devices, err := e.getDeviceList()
	if err != nil {
		fmt.Printf("Error while getting device list: %s\n", err)
		return
	}

	for _, device := range devices.Devices {
		deviceUsage, err := e.client.DeviceAPIUsage(device.DeviceType, device.DeviceId)
		if err != nil {
			fmt.Printf("Error while getting device api usage info: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%+v\n", deviceUsage)
	}
}
