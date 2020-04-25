package main

import (
  "github.com/arcticfoxnv/awair_api"
  "github.com/arcticfoxnv/awair-exporter/awair"
  "github.com/arcticfoxnv/awair-exporter/awair/mock"
  "github.com/prometheus/client_golang/prometheus/testutil"
  "github.com/stretchr/testify/assert"
  "io/ioutil"
  "bytes"
	"testing"
  "time"
)

func TestAwairCollector(t *testing.T) {

  data, err := ioutil.ReadFile("testdata/metrics.txt")
  if err != nil {
    t.Fail()
  }

  expected := bytes.NewReader(data)

  s := mock.NewMockServer()
  defer s.Close()

  client := awair.NewClient(mock.ACCESS_TOKEN, time.Minute, awair_api.SetHTTPClient(s.Client()))
  c := NewAwairCollector(client)

  err = testutil.CollectAndCompare(c, expected)
  assert.Nil(t, err)
}
