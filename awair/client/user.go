package client

import (
	"awair-exporter/structs"
	"fmt"
)

func (c *Client) DeviceAPIUsage(deviceType string, deviceId int) (*structs.DeviceUsage, error) {
	req, err := c.newGetRequest("v1", fmt.Sprintf("users/self/devices/%s/%d/api-usages", deviceType, deviceId))
	if err != nil {
		return nil, err
	}

	data := new(structs.DeviceUsage)
	if err := c.doRequest(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) Devices() (*structs.DeviceList, error) {
	req, err := c.newGetRequest("v1", "users/self/devices")
	if err != nil {
		return nil, err
	}

	data := new(structs.DeviceList)
	if err := c.doRequest(req, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) UserInfo() (*structs.User, error) {
	req, err := c.newGetRequest("v1", "users/self")
	if err != nil {
		return nil, err
	}

	data := new(structs.User)
	if err := c.doRequest(req, data); err != nil {
		return nil, err
	}

	return data, nil
}
