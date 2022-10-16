package floodSensor

import (
	"bytes"
	decoder "decode-sensor-data-sample-go/internal/decoder"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

// ------------------------------------------------------------------------------
// FloodSensorSettingData
// ------------------------------------------------------------------------------
const setting_size int = 166
const schedule_setting_size int = 64

type FloodSensorSettingData struct {
	cableLength         int16
	sendStartWaterLevel float32
	sendInterval        int32
	aliveSetting        AliveType
	scheduleSetting     ScheduleSetting
	reserve             []byte
	fwVersionMajor      int8
	fwVersionMinor      int8
	fwVersionBuild      int8
	hwVersionMajor      int8
	hwVersionMinor      int8
	hwVersionBuild      int8
	battery             int8
	sysStatus           byte
}

func NewSettingFromBase64(sensorDataBase64 string) (*FloodSensorSettingData, error) {
	buf, err := base64.StdEncoding.DecodeString(sensorDataBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode flood sensor setting data: %w", err)
	}
	floodSettingData, err := NewSettingFromBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode flood sensor setting data: %w", err)
	}
	return floodSettingData, nil
}

func NewSettingFromBytes(buf []byte) (*FloodSensorSettingData, error) {
	if len(buf) != setting_size {
		return nil, fmt.Errorf("invalid buf length %d != %d", len(buf), setting_size)
	}
	var data FloodSensorSettingData
	err := data.convertDataFromByte(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bytes to FloodSensorSettingData: %w", err)
	}
	return &data, nil
}

func (f *FloodSensorSettingData) convertDataFromByte(buf []byte) error {
	f.cableLength = decoder.Int16FromByte(buf[0:2])
	f.sendStartWaterLevel = decoder.Float32FromByte(buf[2:6])
	f.sendInterval = decoder.Int32FromByte(buf[6:10])
	f.aliveSetting = aliveType(buf[10])
	s, err := scheduleSetting(buf[11:75], f.aliveSetting)
	if err != nil {
		return fmt.Errorf("failed to convert schedule setting data: %w", err)
	}
	f.scheduleSetting = s
	f.reserve = buf[75:158]
	f.fwVersionMajor = decoder.Int8FromByte(buf[158:159])
	f.fwVersionMinor = decoder.Int8FromByte(buf[159:160])
	f.fwVersionBuild = decoder.Int8FromByte(buf[160:161])
	f.hwVersionMajor = decoder.Int8FromByte(buf[161:162])
	f.hwVersionMinor = decoder.Int8FromByte(buf[162:163])
	f.hwVersionBuild = decoder.Int8FromByte(buf[163:164])
	f.battery = decoder.Int8FromByte(buf[164:165])
	f.sysStatus = buf[165]
	return nil
}

func (f *FloodSensorSettingData) CableLength() int16 {
	return f.cableLength
}

func (f *FloodSensorSettingData) SendStartWaterLevel() float32 {
	return f.sendStartWaterLevel
}

func (f *FloodSensorSettingData) SendInterval() int32 {
	return f.sendInterval
}

func (f *FloodSensorSettingData) ScheduleSetting() ScheduleSetting {
	return f.scheduleSetting
}

func (f *FloodSensorSettingData) AliveSetting() AliveType {
	return f.aliveSetting
}

func (f *FloodSensorSettingData) FwVersion() string {
	return strconv.Itoa(int(f.fwVersionMajor)) + "." + strconv.Itoa(int(f.fwVersionMinor)) + "." + strconv.Itoa(int(f.fwVersionBuild))
}

func (f *FloodSensorSettingData) HwVersion() string {
	return strconv.Itoa(int(f.hwVersionMajor)) + "." + strconv.Itoa(int(f.hwVersionMinor)) + "." + strconv.Itoa(int(f.hwVersionBuild))
}

func (f *FloodSensorSettingData) Battery() int8 {
	return f.battery
}

func (f *FloodSensorSettingData) SysStatus() byte {
	return f.sysStatus
}

func (f *FloodSensorSettingData) ToString() string {
	return fmt.Sprintln("FloodSensorSettingData{") +
		fmt.Sprintln("cableLength:", f.cableLength) +
		fmt.Sprintln("sendStartWaterLevel:", f.sendStartWaterLevel) +
		fmt.Sprintln("sendInterval:", f.sendInterval) +
		fmt.Sprintln("aliveSetting:", f.aliveSetting) +
		fmt.Sprintln("scheduleSetting:") +
		fmt.Sprint(f.scheduleSetting.toString()) +
		fmt.Sprintln("fwVersionMajor:", f.fwVersionMajor) +
		fmt.Sprintln("fwVersionMinor:", f.fwVersionMinor) +
		fmt.Sprintln("fwVersionBuild:", f.fwVersionBuild) +
		fmt.Sprintln("hwVersionMajor:", f.hwVersionMajor) +
		fmt.Sprintln("hwVersionMinor:", f.hwVersionMinor) +
		fmt.Sprintln("hwVersionBuild:", f.hwVersionBuild) +
		fmt.Sprintln("sysStatus:", f.sysStatus)
}

// ------------------------------------------------------------------------------
// AliveType
// ------------------------------------------------------------------------------
type AliveType struct {
	name        string
	value       byte
	description string
}

var (
	Monthly  = AliveType{name: "Monthly", value: 0x00, description: "日時スケジュール"}
	Interval = AliveType{name: "Interval", value: 0x01, description: "インターバル"}
	Daily    = AliveType{name: "Daily", value: 0x02, description: "毎日スケジュール"}
	Off      = AliveType{name: "Off", value: 0x03, description: "OFF"}
	Unknown  = AliveType{name: "Unknown", value: 0xff, description: "不明"}
)

func aliveType(value byte) AliveType {
	switch value {
	case Monthly.value:
		return Monthly
	case Interval.value:
		return Interval
	case Daily.value:
		return Daily
	case Off.value:
		return Off
	default:
		return Unknown
	}
}

// ------------------------------------------------------------------------------
// ScheduleSetting
// ------------------------------------------------------------------------------
type ScheduleSetting struct {
	scheduleType    AliveType
	interval        int32
	dailySchedule   []time.Time
	monthlySchedule []time.Time
}

func scheduleSetting(buf []byte, aliveType AliveType) (ScheduleSetting, error) {
	if len(buf) != schedule_setting_size {
		return ScheduleSetting{}, fmt.Errorf("invalid buf length %d != %d", len(buf), schedule_setting_size)
	}

	switch aliveType {
	case Interval:
		return ScheduleSetting{scheduleType: Interval, interval: decoder.Int32FromByte(buf[0:4])}, nil
	case Monthly:
		return ScheduleSetting{scheduleType: Monthly, monthlySchedule: makeMonthlyScheduleSetting(buf[4:64])}, nil
	case Daily:
		return ScheduleSetting{scheduleType: Daily, dailySchedule: makeDailyScheduleSetting(buf[4:64])}, nil
	default:
		return ScheduleSetting{}, fmt.Errorf("invalid alive type %v", aliveType)
	}
}

func makeMonthlyScheduleSetting(buf []byte) []time.Time {
	const intervalBytes = 3
	var schedules []time.Time
	for i := 0; i < len(buf); i += intervalBytes {
		sub := buf[i : i+intervalBytes]
		if bytes.Equal(sub, []byte{255, 255, 255}) {
			continue
		}
		day := int(decoder.Int8FromByte(sub[0:1]))
		tmpMin := int(decoder.Int16FromByte(sub[1:3]))
		hour := tmpMin / 60
		min := tmpMin % 60
		schedules = append(schedules, time.Date(2022, time.February, day, hour, min, 0, 0, time.Local))
	}
	return schedules
}

func makeDailyScheduleSetting(buf []byte) []time.Time {
	const intervalBytes = 2
	var schedules []time.Time
	for i := 0; i < len(buf); i += intervalBytes {
		sub := buf[i : i+intervalBytes]
		if bytes.Equal(sub, []byte{255, 255}) {
			continue
		}
		tmpMin := int(decoder.Int16FromByte(sub[0:2]))
		hour := tmpMin / 60
		min := tmpMin % 60
		schedules = append(schedules, time.Date(2022, time.February, 1, hour, min, 0, 0, time.Local))
	}
	return schedules
}

func (f *ScheduleSetting) toString() string {
	if f.scheduleType == Interval {
		return fmt.Sprintln("interval: ", f.interval)
	} else if f.scheduleType == Daily {
		var ret string
		for i, s := range f.dailySchedule {
			ret += fmt.Sprintln("index", i, s.Format(time.Kitchen))
		}
		return ret
	} else if f.scheduleType == Monthly {
		var ret string
		for i, s := range f.dailySchedule {
			ret += fmt.Sprintln("index", i, s.Format("02 15:04"))
		}
	}
	return ""
}
