package structs

import "time"

type Device struct {
	DeviceId     int     `json:"deviceId"`
	DeviceType   string  `json:"deviceType"`
	DeviceUUID   string  `json:"deviceUUID"`
	Latitude     float64 `json:"latitude"`
	LocationName string  `json:"locationName"`
	Longitude    float64 `json:"longitude"`
	Name         string  `json:"name"`
	Preference   string  `json:"preference"`
	RoomType     string  `json:"roomType"`
	SpaceType    string  `json:"spaceType"`
	Timezone     string  `json:"timezone"`
}

type DeviceList struct {
	Devices []Device `json:"devices"`
}

type DeviceData struct {
	Indices   []DeviceIndexData     `json:"indices"`
	Score     float64               `json:"score"`
	Sensors   []DeviceSensorReading `json:"sensors"`
	Timestamp time.Time             `json:"timestamp"`
}

type DeviceDataList struct {
	Data []DeviceData `json:"data"`
}

type DeviceIndexData struct {
	Comp  string  `json:"comp"`
	Value float64 `json:"value"`
}

type DeviceSensorReading struct {
	Comp  string  `json:"comp"`
	Value float64 `json:"value"`
}

type DeviceUsage struct {
	Usages []APIUsage `json:"usages"`
}
