package awair

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/arcticfoxnv/awair_api"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

func TestClientSetCacheTTL(t *testing.T) {
	client := NewClient("abc", time.Minute)
	cache := client.clientCache
	client.SetCacheTTL(2 * time.Minute)
	assert.NotEqual(t, cache, client.clientCache)
}

func TestClientGetDeviceAPIUsage(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer abc123", r.Header.Get("Authorization"))

		data, _ := ioutil.ReadFile("testdata/DeviceAPIUsage.json")

		w.Write([]byte(data))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("abc123", time.Minute, awair_api.SetHTTPClient(httpClient))
	data, err := cli.GetDeviceAPIUsage("awair-c", 0)

	assert.Nil(t, err)
	assert.Equal(t, 4, len(data.Usages))
}

func TestClientGetDeviceList(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer abc123", r.Header.Get("Authorization"))
		w.Write([]byte("{}"))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("abc123", time.Minute, awair_api.SetHTTPClient(httpClient))

	_, found := cli.clientCache.Get(DEVICES_KEY)
	assert.False(t, found)

	_, err := cli.GetDeviceList()
	assert.Nil(t, err)

	_, found = cli.clientCache.Get(DEVICES_KEY)
	assert.True(t, found)

	_, err = cli.GetDeviceList()
	assert.Nil(t, err)
}

func TestClientGetLatestData(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer abc123", r.Header.Get("Authorization"))
		w.Write([]byte("{}"))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("abc123", time.Minute, awair_api.SetHTTPClient(httpClient))

	deviceType := "awair-c"
	deviceId := 0
	cacheKey := fmt.Sprintf(DEVICE_LATEST_KEY_FORMAT, deviceType, deviceId)
	_, found := cli.clientCache.Get(cacheKey)
	assert.False(t, found)

	_, err := cli.GetLatestData(deviceType, deviceId)
	assert.Nil(t, err)

	_, found = cli.clientCache.Get(cacheKey)
	assert.True(t, found)

	_, err = cli.GetLatestData(deviceType, deviceId)
	assert.Nil(t, err)
}

func TestClientGetUserInfo(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer abc123", r.Header.Get("Authorization"))

		data, _ := ioutil.ReadFile("testdata/UserInfo.json")

		w.Write([]byte(data))
	})

	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := NewClient("abc123", time.Minute, awair_api.SetHTTPClient(httpClient))
	data, err := cli.GetUserInfo()

	assert.Nil(t, err)
	assert.Equal(t, "Kim", data.LastName)
	assert.Equal(t, "Steve", data.FirstName)
}
