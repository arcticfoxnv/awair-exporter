package main

import (
	"awair-exporter/awair"
	"fmt"
	"net/http"
)

type ExporterHTTP struct {
	client *awair.Client
}

func NewExporterHTTP(client *awair.Client) *ExporterHTTP {
	return &ExporterHTTP{
		client: client,
	}
}

func (e *ExporterHTTP) serveUsage(w http.ResponseWriter, r *http.Request) {
	userInfo, err := e.client.GetUserInfo()
	if err != nil {
		fmt.Printf("Error while getting user info: %s\n", err)
		return
	}
	fmt.Fprintf(w, "%+v\n", userInfo)

	devices, err := e.client.GetDeviceList()
	if err != nil {
		fmt.Printf("Error while getting device list: %s\n", err)
		return
	}

	for _, device := range devices.Devices {
		deviceUsage, err := e.client.GetDeviceAPIUsage(device.DeviceType, device.DeviceId)
		if err != nil {
			fmt.Printf("Error while getting device api usage info: %s\n", err)
			return
		}
		fmt.Fprintf(w, "%+v\n", deviceUsage)
	}
}
