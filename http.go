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

const (
	DEVICES_KEY              = "devices"
	DEVICE_LATEST_KEY_FORMAT = "device-latest-%s-%d"
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

func (e *ExporterHTTP) getLatestData(deviceType string, deviceId int) (*structs.DeviceDataList, error) {
	cacheKey := fmt.Sprintf(DEVICE_LATEST_KEY_FORMAT, deviceType, deviceId)
	if data, found := e.clientCache.Get(cacheKey); found {
		return data.(*structs.DeviceDataList), nil
	}
	log.Printf("Fetching data for %s-%d", deviceType, deviceId)

	data, err := e.client.UserLatestAirData(deviceType, deviceId)
	if err != nil {
		return nil, err
	}

	e.clientCache.Set(cacheKey, data, cache.DefaultExpiration)
	return data, nil
}

func (e *ExporterHTTP) serveLatest(w http.ResponseWriter, r *http.Request) {
	devices, err := e.getDeviceList()
	if err != nil {
		fmt.Printf("Error while getting device list: %s\n", err)
		return
	}

	for _, device := range devices.Devices {
		dataList, err := e.getLatestData(device.DeviceType, device.DeviceId)
		if err != nil {
			fmt.Printf("Error while getting latest air data: %s\n", err)
			return
		}

		for _, data := range dataList.Data {
			fmt.Fprintf(w, "%s\n", formatDeviceMetrics(&device, &data))
		}
	}
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
