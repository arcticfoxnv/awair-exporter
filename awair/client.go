package awair

import "awair-exporter/awair/client"

func NewClient(accessToken string) *client.Client {
	return &client.Client{
		AccessToken: accessToken,
	}
}
