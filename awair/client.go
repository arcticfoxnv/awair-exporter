package awair

import (
	"fmt"
	"github.com/arcticfoxnv/awair_api"
	"github.com/patrickmn/go-cache"
	"log"
	"sync"
	"time"
)

const (
	DEVICES_KEY              = "devices"
	DEVICE_LATEST_KEY_FORMAT = "device-latest-%s-%d"
)

type Client struct {
	client *awair_api.Client

	clientCacheLock *sync.RWMutex
	clientCache     *cache.Cache
}

func NewClient(accessToken string, cacheTTL time.Duration, options ...awair_api.Option) *Client {
	cli := &Client{
		client:          awair_api.NewClient(accessToken, options...),
		clientCacheLock: new(sync.RWMutex),
	}

	cli.SetCacheTTL(cacheTTL)

	return cli
}

func (c *Client) SetCacheTTL(ttl time.Duration) {
	c.clientCacheLock.Lock()
	defer c.clientCacheLock.Unlock()

	if c.clientCache != nil {
		c.clientCache.Flush()
	}

	c.clientCache = cache.New(ttl, 10*time.Minute)
}

func (c *Client) GetDeviceAPIUsage(deviceType string, deviceId int) (*awair_api.DeviceUsage, error) {
	return c.client.DeviceAPIUsage(deviceType, deviceId)
}

func (c *Client) GetDeviceList() (*awair_api.DeviceList, error) {
	c.clientCacheLock.RLock()
	defer c.clientCacheLock.RUnlock()

	if data, found := c.clientCache.Get(DEVICES_KEY); found {
		return data.(*awair_api.DeviceList), nil
	}
	log.Printf("Fetching device list")

	devices, err := c.client.Devices()
	if err != nil {
		return nil, err
	}

	c.clientCache.Set(DEVICES_KEY, devices, cache.DefaultExpiration)
	return devices, nil
}

func (c *Client) GetLatestData(deviceType string, deviceId int) (*awair_api.DeviceDataList, error) {
	c.clientCacheLock.RLock()
	defer c.clientCacheLock.RUnlock()

	cacheKey := fmt.Sprintf(DEVICE_LATEST_KEY_FORMAT, deviceType, deviceId)
	if data, found := c.clientCache.Get(cacheKey); found {
		return data.(*awair_api.DeviceDataList), nil
	}
	log.Printf("Fetching data for %s-%d", deviceType, deviceId)

	data, err := c.client.UserLatestAirData(deviceType, deviceId)
	if err != nil {
		return nil, err
	}

	c.clientCache.Set(cacheKey, data, cache.DefaultExpiration)
	return data, nil
}

func (c *Client) GetUserInfo() (*awair_api.User, error) {
	return c.client.UserInfo()
}
