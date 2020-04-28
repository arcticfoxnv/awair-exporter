package main

import (
	"github.com/arcticfoxnv/awair-exporter/awair"
	"github.com/arcticfoxnv/awair-exporter/awair/mock"
	"github.com/arcticfoxnv/awair_api"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetMetaUsage(t *testing.T) {
	req, err := http.NewRequest("GET", "/meta/usage", nil)
	if err != nil {
		t.Fail()
	}

	s := mock.NewMockServer()
	defer s.Close()

	cli := awair.NewClient(mock.ACCESS_TOKEN, time.Minute, awair_api.SetHTTPClient(s.Client()))

	rr := httptest.NewRecorder()
	exporter := NewExporter(cli)
	exporter.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d", status, http.StatusOK)
	}

	testString := "Email:developer@getawair.com"
	if !strings.Contains(rr.Body.String(), testString) {
		t.Errorf("handler returned unexpected body: got %s, want it to contain %s",
			rr.Body.String(), testString)
	}
}

func TestGetMetrics(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fail()
	}

	s := mock.NewMockServer()
	defer s.Close()

	cli := awair.NewClient(mock.ACCESS_TOKEN, time.Minute, awair_api.SetHTTPClient(s.Client()))

	rr := httptest.NewRecorder()
	exporter := NewExporter(cli)
	exporter.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d want %d", status, http.StatusOK)
	}

	testString := `awair_score{device_name="my awair",device_type="awair",device_uuid="awair_0",location_name="my home",room_type="bedroom",space_type="home"} 90`
	if !strings.Contains(rr.Body.String(), testString) {
		t.Errorf("handler returned unexpected body: got %s, want it to contain %s",
			rr.Body.String(), testString)
	}
}
