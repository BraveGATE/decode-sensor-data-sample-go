package main

import (
	"bytes"
	"fmt"
	"time"
)

type Webhook struct {
	Application Application `json:"application"`
	Router      Router      `json:"router"`
	Device      Device      `json:"device"`
	UplinkId    string      `json:"uplink_id"`
	Date        JsonTime    `json:"date"`
}

type Application struct {
	ApplicationId string `json:"application_id"`
	Name          string `json:"name"`
}

type Router struct {
	RouterId  string `json:"router_id"`
	Imsi      string `json:"imsi"`
	Rssi      int    `json:"rssi"`
	Battery   int    `json:"battery"`
	FwVersion string `json:"fw_version"`
}

type Device struct {
	DeviceId   string                 `json:"device_id"`
	SensorId   string                 `json:"sensor_id"`
	SensorName string                 `json:"sensor_name"`
	Rssi       int                    `json:"rssi"`
	Data       map[string]interface{} `json:"data"`
}

// JsonTime exists so that we can have a String method converting the date
type JsonTime string

// String converts the unix timestamp into a string
func (t *JsonTime) String() string {
	tm := t.Time()
	return fmt.Sprintf("\"%s\"", tm.Format(time.RFC3339))
}

// Time returns a `time.Time` representation of this value.
func (t *JsonTime) Time() time.Time {
	tt, _ := time.Parse(time.RFC3339, string(*t))
	return tt
}

// UnmarshalJSON will unmarshal both string and int JSON values
func (t *JsonTime) UnmarshalJSON(buf []byte) error {
	s := bytes.Trim(buf, `"`)
	*t = JsonTime(string(s))
	return nil
}
