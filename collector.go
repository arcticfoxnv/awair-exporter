package main

import (
	"fmt"
	"github.com/arcticfoxnv/awair_api"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	DEVICES_KEY              = "devices"
	DEVICE_LATEST_KEY_FORMAT = "device-latest-%s-%d"
)

var collectorLabels = []string{
	"device_name",
	"device_type",
	"device_uuid",
	"location_name",
	"room_type",
	"space_type",
}

var awairScoreGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "awair",
	Name:      "score",
	Help:      "Awair Score Value",
}, collectorLabels)

type AwairCollector struct {
	collectLock *sync.Mutex

	client      *awair_api.Client
	clientCache *cache.Cache
}

func NewAwairCollector(client *awair_api.Client, cacheTTL time.Duration) *AwairCollector {
	return &AwairCollector{
		client:      client,
		clientCache: cache.New(cacheTTL, 10*time.Minute),
		collectLock: new(sync.Mutex),
	}
}

func (ac *AwairCollector) Describe(ch chan<- *prometheus.Desc) {
	awairScoreGauge.Describe(ch)
}

func (ac *AwairCollector) Collect(ch chan<- prometheus.Metric) {
	ac.collectLock.Lock()
	defer ac.collectLock.Unlock()

	devices, err := ac.getDeviceList()
	if err != nil {
		fmt.Printf("Error while getting device list: %s\n", err)
		return
	}

	for _, device := range devices.Devices {
		dataList, err := ac.getLatestData(device.DeviceType, device.DeviceId)
		if err != nil {
			fmt.Printf("Error while getting latest air data: %s\n", err)
			return
		}

		labels := make(prometheus.Labels)
		labels["device_name"] = device.Name
		labels["device_type"] = device.DeviceType
		labels["device_uuid"] = device.DeviceUUID
		labels["location_name"] = device.LocationName
		labels["room_type"] = device.RoomType
		labels["space_type"] = device.SpaceType

		for k, v := range labels {
			labels[k] = strings.ToLower(v)
		}

		for _, data := range dataList.Data {
			// Awair Composite Score
			awairScoreGauge.With(labels).Set(data.Score)

			// Sensor values
			for _, sensorData := range data.Sensors {
				gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "awair",
					Subsystem: "sensor",
					Name:      strings.ToLower(sensorData.Comp),
				}, collectorLabels)
				gauge.With(labels).Set(sensorData.Value)
				gauge.Collect(ch)
			}

			// Index values
			for _, sensorData := range data.Indices {
				gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: "awair",
					Subsystem: "index",
					Name:      strings.ToLower(sensorData.Comp),
				}, collectorLabels)
				gauge.With(labels).Set(sensorData.Value)
				gauge.Collect(ch)
			}
		}
	}

	awairScoreGauge.Collect(ch)
}

func (ac *AwairCollector) getDeviceList() (*awair_api.DeviceList, error) {
	if data, found := ac.clientCache.Get(DEVICES_KEY); found {
		return data.(*awair_api.DeviceList), nil
	}
	log.Printf("Fetching device list")

	devices, err := ac.client.Devices()
	if err != nil {
		return nil, err
	}

	ac.clientCache.Set(DEVICES_KEY, devices, cache.DefaultExpiration)
	return devices, nil
}

func (ac *AwairCollector) getLatestData(deviceType string, deviceId int) (*awair_api.DeviceDataList, error) {
	cacheKey := fmt.Sprintf(DEVICE_LATEST_KEY_FORMAT, deviceType, deviceId)
	if data, found := ac.clientCache.Get(cacheKey); found {
		return data.(*awair_api.DeviceDataList), nil
	}
	log.Printf("Fetching data for %s-%d", deviceType, deviceId)

	data, err := ac.client.UserLatestAirData(deviceType, deviceId)
	if err != nil {
		return nil, err
	}

	ac.clientCache.Set(cacheKey, data, cache.DefaultExpiration)
	return data, nil
}
