package main

import (
  "fmt"
  "github.com/arcticfoxnv/awair-exporter/awair"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "net/http"
)

type Exporter struct {
  client  *awair.Client
  httpMux *http.ServeMux
  registry *prometheus.Registry
}

func NewExporter(client *awair.Client) *Exporter {
  e := &Exporter{
    client: client,
    httpMux: http.NewServeMux(),
    registry: prometheus.NewRegistry(),
  }

  e.registry.MustRegister(NewAwairCollector(client))

  e.httpMux.Handle("/metrics", promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{}))
  e.httpMux.HandleFunc("/meta/usage", e.handleMetaUsage)

  return e
}

func (e *Exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  e.httpMux.ServeHTTP(w, r)
}

func (e *Exporter) handleMetaUsage(w http.ResponseWriter, r *http.Request) {
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
