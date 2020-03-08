package client

import (
	"awair-exporter/awair/structs"
	"fmt"
)

func (c *Client) UserLatestAirData(deviceType string, deviceId int) (*structs.DeviceDataList, error) {
	endpoint := fmt.Sprintf("users/self/devices/%s/%d/air-data/latest", deviceType, deviceId)
	req, err := c.newGetRequest("v1", endpoint)
	if err != nil {
		return nil, err
	}

	if c.UseFarenheit {
		c.appendQueryParam(req, "farenheit", "true")
	}

	data := new(structs.DeviceDataList)
	if err := c.doRequest(req, data); err != nil {
		return nil, err
	}

	return data, nil
}
