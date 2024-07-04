package main

import (
	"encoding/json"
	"time"
)

type TempstickResponse struct {
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SensorData struct {
	ID                    string        `json:"id"`
	Version               string        `json:"version"`
	SensorID              string        `json:"sensor_id"`
	SensorName            string        `json:"sensor_name"`
	SensorMacAddr         string        `json:"sensor_mac_addr"`
	OwnerID               string        `json:"owner_id"`
	Type                  string        `json:"type"`
	AlertInterval         string        `json:"alert_interval"`
	SendInterval          string        `json:"send_interval"`
	LastTemp              float64       `json:"last_temp"`
	LastHumidity          float64       `json:"last_humidity"`
	LastVoltage           float64       `json:"last_voltage"`
	WifiConnectTime       int           `json:"wifi_connect_time"`
	RSSI                  int           `json:"rssi"`
	LastCheckin           time.Time     `json:"last_checkin"`
	NextCheckin           time.Time     `json:"next_checkin"`
	SSID                  string        `json:"ssid"`
	Offline               string        `json:"offline"`
	Alerts                []interface{} `json:"alerts"`
	UseSensorSettings     int           `json:"use_sensor_settings"`
	TempOffset            int           `json:"temp_offset"`
	HumidityOffset        int           `json:"humidity_offset"`
	AlertTempBelow        string        `json:"alert_temp_below"`
	AlertTempAbove        string        `json:"alert_temp_above"`
	AlertHumidityBelow    string        `json:"alert_humidity_below"`
	AlertHumidityAbove    string        `json:"alert_humidity_above"`
	ConnectionSensitivity int           `json:"connection_sensitivity"`
	HIM                   int           `json:"HI_M"`
	HI                    int           `json:"HI"`
	DPM                   int           `json:"DP_M"`
	DP                    int           `json:"DP"`
	WlanA                 string        `json:"wlanA"`
	WlanB                 string        `json:"wlanB"`
	LastWlan              string        `json:"last_wlan"`
	Wlan1Used             time.Time     `json:"wlan_1_used"`
	UseAlertInterval      int           `json:"use_alert_interval"`
	UseOffset             int           `json:"use_offset"`
	BatteryPct            int           `json:"battery_pct"`
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
	if len(s) > 0 && s[len(s)-1] == 'Z' {
		return s[:len(s)-1]
	}
	return s
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
		var sensorData struct {
			SensorData
			LastCheckin string `json:"last_checkin"`
			NextCheckin string `json:"next_checkin"`
			Wlan1Used   string `json:"wlan_1_used"`
		}
		if err := json.Unmarshal(aux.Data, &sensorData); err != nil {
			return err
		}
		parsedLastCheckin, err := time.Parse("2006-01-02 15:04:05-07:00", removeTrailingZ(sensorData.LastCheckin))
		if err != nil {
			return err
		}
		parsedNextCheckin, err := time.Parse("2006-01-02 15:04:05-07:00", removeTrailingZ(sensorData.NextCheckin))
		if err != nil {
			return err
		}
		parsedWlan1Used, err := time.Parse("2006-01-02 15:04:05", removeTrailingZ(sensorData.Wlan1Used))
		if err != nil {
			return err
		}
		sensorData.SensorData.LastCheckin = parsedLastCheckin
		sensorData.SensorData.NextCheckin = parsedNextCheckin
		sensorData.SensorData.Wlan1Used = parsedWlan1Used
		r.Data = &sensorData.SensorData
	default:
		r.Data = aux.Data
	}
	return nil
}
