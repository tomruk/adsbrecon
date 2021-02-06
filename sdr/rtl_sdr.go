package sdr

import (
	"fmt"
	"strconv"
	"sync"

	rtl "github.com/jpoirier/gortlsdr"
)

type RTLSDRDevice struct {
	ID         int
	Index      int
	Name       string
	USBStrings string

	dev *rtl.Context

	runMutex *sync.Mutex
}

func (device RTLSDRDevice) Run() error {
	device.RunLock()
	defer device.RunUnlock()

	var err error
	device.dev, err = rtl.Open(device.Index)
	if err != nil {
		return err
	}

	device.dev.SetAgcMode(true)
	device.dev.SetSampleRate(2000000)
	device.dev.SetTunerGainMode(false)
	err = device.dev.SetCenterFreq(1090000000)
	if err != nil {
		return err
	}

	device.dev.ResetBuffer()

	fmt.Printf("Starting capture with RTL-SDR (index: #%d)", device.Index)

	go func() {
		readAsyncErr := device.dev.ReadAsync(rtlsdrCallback, nil, 0, DATA_LEN)
		if readAsyncErr != nil {
			fmt.Println()
			fmt.Printf("ReadAsync error: %s", readAsyncErr.Error())
		}

		noDevicePrompt()
	}()

	return nil
}

func rtlsdrCallback(buf []byte) {
	callback(buf)
}

func (device RTLSDRDevice) Stop() error {
	device.RunLock()
	defer device.RunUnlock()

	fmt.Println("Stopping the RTL-SDR device")

	if device.dev != nil {
		device.dev.CancelAsync()
		return device.dev.Close()
	}

	return nil
}

func (device RTLSDRDevice) RunLock() {
	device.runMutex.Lock()
}

func (device RTLSDRDevice) RunUnlock() {
	device.runMutex.Unlock()
}

func (device RTLSDRDevice) GetDeviceType() string {
	return "RTL-SDR"
}

func (device RTLSDRDevice) GetID() int {
	return device.ID
}

func (device RTLSDRDevice) GetIndex() int {
	return device.Index
}

func (device RTLSDRDevice) GetPrintableBox() []string {
	return []string{
		strconv.Itoa(device.ID),
		strconv.Itoa(device.Index),
		device.Name,
		device.USBStrings,
	}
}

func GetRTLSDRDevices(offset int) (devices []SDRDevice) {
	deviceCount := rtl.GetDeviceCount()

	for index := 0; index < deviceCount; index, offset = index+1, offset+1 {
		dev := NewRTLSDRDevice(index, offset)
		devices = append(devices, dev)
	}

	return
}

func NewRTLSDRDevice(id, index int) RTLSDRDevice {
	dev := RTLSDRDevice{
		ID:       id,
		Index:    index,
		Name:     rtl.GetDeviceName(index),
		runMutex: &sync.Mutex{},
	}

	manufacturer, product, serialnumber, err := rtl.GetDeviceUsbStrings(index)
	if err == nil {
		dev.USBStrings = fmt.Sprintf("%s/%s/%s", manufacturer, product, serialnumber)
	} else {
		dev.USBStrings = "N/A"
	}

	return dev
}
