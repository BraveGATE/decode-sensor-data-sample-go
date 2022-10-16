package floodSensor

import (
	decoder "decode-sensor-data-sample-go/internal/decoder"
	"encoding/base64"
	"fmt"
	"strconv"
)

// ------------------------------------------------------------------------------
// FloodSensorData
// ------------------------------------------------------------------------------
const data_size int = 20

type FloodSensorData struct {
	fwVersionMajor   int8
	fwVersionMinor   int8
	fwVersionBuild   int8
	waterPressure    float32
	waterTemperature float32
	airPressure      float32
	airTemperature   float32
	battery          int8
}

func NewDataFromBase64(sensorDataBase64 string) (*FloodSensorData, error) {
	buf, err := base64.StdEncoding.DecodeString(sensorDataBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode flood sensor data: %w", err)
	}
	floodSensorData, err := NewDataFromBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode flood sensor data: %w", err)
	}
	return floodSensorData, nil
}

func NewDataFromBytes(buf []byte) (*FloodSensorData, error) {
	if len(buf) != data_size {
		return nil, fmt.Errorf("invalid buf length %d != %d", len(buf), data_size)
	}
	var data FloodSensorData
	data.convertDataFromByte(buf)
	return &data, nil
}

func (f *FloodSensorData) convertDataFromByte(buf []byte) {
	f.fwVersionMajor = decoder.Int8FromByte(buf[0:1])
	f.fwVersionMinor = decoder.Int8FromByte(buf[1:2])
	f.fwVersionBuild = decoder.Int8FromByte(buf[2:3])
	f.waterPressure = decoder.Float32FromByte(buf[3:7])
	f.waterTemperature = decoder.Float32FromByte(buf[7:11])
	f.airPressure = decoder.Float32FromByte(buf[11:15])
	f.airTemperature = decoder.Float32FromByte(buf[15:19])
	f.battery = decoder.Int8FromByte(buf[19:])
}

func (f *FloodSensorData) FwVersion() string {
	return strconv.Itoa(int(f.fwVersionMajor)) + "." + strconv.Itoa(int(f.fwVersionMinor)) + "." + strconv.Itoa(int(f.fwVersionBuild))
}

func (f *FloodSensorData) WaterPressure() float32 {
	return f.waterPressure
}

func (f *FloodSensorData) WaterTemperature() float32 {
	return f.waterTemperature
}

func (f *FloodSensorData) AirPressure() float32 {
	return f.airPressure
}

func (f *FloodSensorData) AirTemperature() float32 {
	return f.airTemperature
}

func (f *FloodSensorData) Battery() int8 {
	return f.battery
}

func (f *FloodSensorData) ToString() string {
	return fmt.Sprintf("FloodSensorData{fwVersion: %s, waterPressure: %v, waterTemperature: %v, airPressure: %v, airTemperature: %v, battery: %v}",
		f.FwVersion(), f.WaterPressure(), f.WaterTemperature(), f.AirPressure(), f.AirTemperature(), f.Battery())
}
