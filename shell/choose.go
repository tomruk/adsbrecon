package shell

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tomruk/adsbrecon/sdr"
)

func chooseDevice(confirm bool) (device sdr.SDRDevice) {
	device = nil
	devices := sdr.GetAllSDRDevices()
	devicesLen := len(devices)

	if devicesLen == 0 {
		fmt.Println("No SDR device found")
		return
	}

	var data [][]string

	for _, d := range devices {
		data = append(data, d.GetPrintableBox())
	}

	fmt.Println()

	if devicesLen == 1 {
		fmt.Println("Found 1 device")
	} else {
		fmt.Printf("Found %d devices\n", devicesLen)
	}

	printTable([]string{
		"ID",
		"INDEX",
		"NAME",
		"DETAILS",
	}, nil, data, true)

	if devicesLen == 1 {
		if confirm {
			if confirmation(fmt.Sprintf("\nOnly 1 device found, shall we continue to use this device?"), true) == false {
				return
			}
		}

		device = devices[0]
		return
	}

	shell.Print("Device id to use: ")
	idInput := strings.TrimSpace(shell.ReadLine())
	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Println("Please input a positive integer")
		return
	}

	for _, d := range devices {
		if d.GetID() == id {
			device = d
		}
	}

	if device == nil {
		fmt.Printf("No device with id #%d found\n", id)
	}

	return
}
