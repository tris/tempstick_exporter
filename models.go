package main

import (
	"encoding/json"
	"strings"
	"time"
)

const checkinTimeLayout = "2006-01-02 15:04:05-07:00"

type CheckinTime struct {
	time.Time
}

type TempstickResponse struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SensorData struct {
	ID                    json.Number   `json:"id"`
	Version               string        `json:"version"`
	SensorID              string        `json:"sensor_id"`
	SensorName            string        `json:"sensor_name"`
	SensorMacAddr         string        `json:"sensor_mac_addr"`
	OwnerID               string        `json:"owner_id"`
	Type                  string        `json:"type"`
	AlertInterval         json.Number   `json:"alert_interval"`
	SendInterval          json.Number   `json:"send_interval"`
	LastTemp              float64       `json:"last_temp"`
	LastHumidity          float64       `json:"last_humidity"`
	LastVoltage           float64       `json:"last_voltage"`
	WifiConnectTime       json.Number   `json:"wifi_connect_time"`
	RSSI                  float64       `json:"rssi"`
	LastCheckin           CheckinTime   `json:"last_checkin"`
	NextCheckin           CheckinTime   `json:"next_checkin"`
	SSID                  string        `json:"ssid"`
	Offline               json.Number   `json:"offline"`
	Alerts                []interface{} `json:"alerts"`
	UseSensorSettings     json.Number   `json:"use_sensor_settings"`
	TempOffset            json.Number   `json:"temp_offset"`
	HumidityOffset        json.Number   `json:"humidity_offset"`
	AlertTempBelow        string        `json:"alert_temp_below"`
	AlertTempAbove        string        `json:"alert_temp_above"`
	AlertHumidityBelow    string        `json:"alert_humidity_below"`
	AlertHumidityAbove    string        `json:"alert_humidity_above"`
	ConnectionSensitivity json.Number   `json:"connection_sensitivity"`
	HIM                   json.Number   `json:"HI_M"`
	HI                    json.Number   `json:"HI"`
	DPM                   json.Number   `json:"DP_M"`
	DP                    json.Number   `json:"DP"`
	WlanA                 string        `json:"wlanA"`
	WlanB                 string        `json:"wlanB"`
	LastWlan              string        `json:"last_wlan"`
	UseAlertInterval      json.Number   `json:"use_alert_interval"`
	UseOffset             json.Number   `json:"use_offset"`
	BatteryPct            json.Number   `json:"battery_pct"`
	LastMessages          []LastMessage `json:"last_messages"`
}

type LastMessage struct {
	Temperature   float64 `json:"temperature"`
	Humidity      float64 `json:"humidity"`
	Voltage       string  `json:"voltage"`
	RSSI          string  `json:"RSSI"`
	TimeToConnect string  `json:"time_to_connect"`
	SensorTimeUTC string  `json:"sensor_time_utc"`
}

func removeTrailingZ(s string) string {
	s = strings.Trim(s, `"`)
	if len(s) > 0 && s[len(s)-1] == 'Z' {
		return s[:len(s)-1]
	}
	return s
}

func (ct *CheckinTime) UnmarshalJSON(data []byte) error {
	strInput := string(data)
	strInput = removeTrailingZ(strInput)
	parsedTime, err := time.Parse(checkinTimeLayout, strInput)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

func (r *TempstickResponse) UnmarshalJSON(data []byte) error {
	type Alias TempstickResponse
	aux := &struct {
		Data json.RawMessage `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch r.Message {
	case "get sensor":
		var sensorData SensorData
		if err := json.Unmarshal(aux.Data, &sensorData); err != nil {
			return err
		}
		r.Data = &sensorData
	default:
		r.Data = aux.Data
	}
	return nil
}
