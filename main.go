package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	FloodSensor "github.com/BraveGATE/decode-sensor-data-sample-go/sensor/floodsensor"
)

func main() {

	// ----------------------------------------------------
	// Parse FloodSensorData
	// ----------------------------------------------------
	base64string, err := GetSensorDataFromWebhookJson("./testdata/webhook/floodsensor/sensor_data.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	floodSensorData, err := FloodSensor.NewDataFromBase64(base64string)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(floodSensorData.ToString())

	// ----------------------------------------------------
	// Parse FloodSensorSettingData
	// ----------------------------------------------------
	base64string, err = GetSensorDataFromWebhookJson("./testdata/webhook/floodsensor/sensor_setting.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	floodSensorSettingData, err := FloodSensor.NewSettingFromBase64(base64string)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(floodSensorSettingData.ToString())
}

func GetSensorDataFromWebhookJson(jsonPath string) (string, error) {
	path := filepath.Join(jsonPath)
	jsonText, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("could not read webhook json file: %v", err)
	}

	var webhook Webhook
	json.Unmarshal([]byte(jsonText), &webhook)
	data := webhook.Device.Data["data"]
	base64string, ok := data.(string)
	if !ok {
		return "", fmt.Errorf("could not parse webhook data: %v", webhook)
	}
	return base64string, nil
}
