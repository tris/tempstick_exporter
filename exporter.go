package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiKey = ""

	labelNames = []string{
		"latitude",
		"longitude",
	}
)

func init() {
	apiKey = os.Getenv("TEMPSTICK_API_KEY")
	if apiKey == "" {
		log.Fatal("TEMPSTICK_API_KEY environment variable is required")
	}
}

func scrapeHandler(w http.ResponseWriter, r *http.Request) {
	sensorID := r.URL.Query().Get("sensor_id")
	if sensorID == "" {
		http.Error(w, "Missing required parameter: sensor_id", http.StatusBadRequest)
		return
	}

	// Fetch current readings
	// API docs: https://tempstickapi.com/docs/
	req, err := http.NewRequest("GET", "https://tempstickapi.com/api/v1/sensor/"+sensorID, nil)
	req.Header.Add("X-API-KEY", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response TempstickResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Println(err)
		return
	}

	// Check for success
	if response.Type != "success" {
		log.Println("Temp Stick API returned an error")
		return
	}

	var sensorData *SensorData
	if data, ok := response.Data.(*SensorData); ok {
                sensorData = data
	} else {
		log.Println("Failed to assert SensorData")
		return
	}

	// Create a new registry for this scrape
	registry := prometheus.NewRegistry()

	// Create labels
	labels := prometheus.Labels{
		"sensor_id":       sensorData.SensorID,
		"sensor_name":     sensorData.SensorName,
		"sensor_mac_addr": sensorData.SensorMacAddr,
	}

	// Create gauges and set values
	collector := &TempstickCollector{
		metrics: []*tempstickMetric{
			newTempstickMetric("tempstick_temp", "Temperature in degrees Celsius", labels, sensorData.LastTemp, sensorData.LastCheckin.Time),
			newTempstickMetric("tempstick_humidity", "Humidity in percent", labels, sensorData.LastHumidity, sensorData.LastCheckin.Time),
			newTempstickMetric("tempstick_voltage", "Battery voltage", labels, sensorData.LastVoltage, sensorData.LastCheckin.Time),
			newTempstickMetric("tempstick_rssi", "RSSI in dBm", labels, sensorData.RSSI, sensorData.LastCheckin.Time),
		},
	}

	// Register the gauges with the registry
	registry.MustRegister(collector)

	// Use a promhttp.HandlerFor with the new registry to serve the metrics
	promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
