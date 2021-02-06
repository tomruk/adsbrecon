package sdr

import (
	"fmt"
)

var (
	Device SDRDevice
)

type SDRDevice interface {
	Run() error
	Stop() error

	RunLock()
	RunUnlock()

	GetDeviceType() string
	GetID() int
	GetIndex() int
	GetPrintableBox() []string
}

func GetAllSDRDevices() (devices []SDRDevice) {
	offset := 0

	rtlsdrDevices := GetRTLSDRDevices(offset)
	offset += len(rtlsdrDevices)

	devices = append(devices, rtlsdrDevices...)
	return devices
}

func noDevicePrompt() {
	fmt.Println("\nNo device is receiving right now, choose one with 'choose-device' command")
}
