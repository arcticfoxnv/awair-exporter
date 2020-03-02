package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	AwairV1 = "https://developer-apis.awair.is/v1"
)

type Client struct {
	AccessToken  string
	UseFarenheit bool
}

func (c *Client) getEndpoint(version, endpoint string) string {
	var base string

	switch version {
	case "v1":
		base = AwairV1
	default:
		base = AwairV1
	}

	return fmt.Sprintf("%s/%s", base, endpoint)
}

func (c *Client) appendQueryParam(req *http.Request, k, v string) {
	q := req.URL.Query()
	q.Set(k, v)
	req.URL.RawQuery = q.Encode()
}

func (c *Client) newGetRequest(version, endpoint string) (*http.Request, error) {
	url := c.getEndpoint(version, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	return req, nil
}

func (c *Client) doRequest(req *http.Request, data interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return err
	}

	return nil
}
