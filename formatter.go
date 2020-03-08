package main

import (
	"awair-exporter/awair/structs"
	"fmt"
	"sort"
	"strings"
)

const (
	AWAIR_SCORE        = "awair_score {%s} %f"
	AWAIR_SENSOR_VALUE = "awair_sensor_%s {%s} %f"
	AWAIR_INDEX_VALUE  = "awair_index_%s {%s} %f"
	STRING_LABEL       = "%s=\"%s\""
)

func formatDeviceMetrics(device *structs.Device, deviceData *structs.DeviceData) string {
	labels := strings.ToLower(strings.Join([]string{
		fmt.Sprintf(STRING_LABEL, "device_name", device.Name),
		fmt.Sprintf(STRING_LABEL, "device_type", device.DeviceType),
		fmt.Sprintf(STRING_LABEL, "device_uuid", device.DeviceUUID),
		fmt.Sprintf(STRING_LABEL, "location_name", device.LocationName),
		fmt.Sprintf(STRING_LABEL, "room_type", device.RoomType),
		fmt.Sprintf(STRING_LABEL, "space_type", device.SpaceType),
	}, ","))
	output := make([]string, 0)

	// Awair Composite Score
	output = append(output, fmt.Sprintf(AWAIR_SCORE, labels, deviceData.Score))

	// Sensor values
	sensorOutput := make([]string, 0)
	for _, sensorData := range deviceData.Sensors {
		compName := strings.ToLower(sensorData.Comp)
		sensorOutput = append(sensorOutput, fmt.Sprintf(AWAIR_SENSOR_VALUE, compName, labels, sensorData.Value))
	}
	sort.Strings(sensorOutput)
	output = append(output, sensorOutput...)

	// Index values
	indicesOutput := make([]string, 0)
	for _, sensorData := range deviceData.Indices {
		compName := strings.ToLower(sensorData.Comp)
		indicesOutput = append(indicesOutput, fmt.Sprintf(AWAIR_INDEX_VALUE, compName, labels, sensorData.Value))
	}
	sort.Strings(indicesOutput)
	output = append(output, indicesOutput...)

	return strings.Join(output, "\n")
}
