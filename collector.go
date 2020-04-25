package main

import (
	"fmt"
	"github.com/arcticfoxnv/awair-exporter/awair"
	"github.com/arcticfoxnv/awair_api"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
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
	client      *awair.Client
	collectLock *sync.Mutex
}

func NewAwairCollector(client *awair.Client) *AwairCollector {
	return &AwairCollector{
		client:      client,
		collectLock: new(sync.Mutex),
	}
}

func (ac *AwairCollector) Describe(ch chan<- *prometheus.Desc) {
	awairScoreGauge.Describe(ch)
}

func (ac *AwairCollector) Collect(ch chan<- prometheus.Metric) {
	ac.collectLock.Lock()
	defer ac.collectLock.Unlock()

	devices, err := ac.client.GetDeviceList()
	if err != nil {
		fmt.Printf("Error while getting device list: %s\n", err)
		return
	}

	for _, device := range devices.Devices {
		dataList, err := ac.client.GetLatestData(device.DeviceType, device.DeviceId)
		if err != nil {
			fmt.Printf("Error while getting latest air data: %s\n", err)
			continue
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
			ac.collectSensors(ch, data.Sensors, &labels)

			// Index values
			ac.collectIndices(ch, data.Indices, &labels)
		}
	}

	awairScoreGauge.Collect(ch)
}

func (ac *AwairCollector) collectSensors(ch chan<- prometheus.Metric, sensors []awair_api.DeviceSensorReading, labels *prometheus.Labels) {
	for _, sensorData := range sensors {
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "awair",
			Subsystem: "sensor",
			Name:      strings.ToLower(sensorData.Comp),
		}, collectorLabels)
		gauge.With(*labels).Set(sensorData.Value)
		gauge.Collect(ch)
	}
}

func (ac *AwairCollector) collectIndices(ch chan<- prometheus.Metric, indices []awair_api.DeviceIndexData, labels *prometheus.Labels) {
	for _, sensorData := range indices {
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "awair",
			Subsystem: "index",
			Name:      strings.ToLower(sensorData.Comp),
		}, collectorLabels)
		gauge.With(*labels).Set(sensorData.Value)
		gauge.Collect(ch)
	}
}
