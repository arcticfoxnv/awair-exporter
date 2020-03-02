package structs

type APIQuota struct {
	Quota int    `json:"quota"`
	Scope string `json:"scope"`
}

type APIUsage struct {
	Scope string `json:"scope"`
	Usage int    `json:"usage"`
}
