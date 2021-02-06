package shell

import (
	"fmt"

	"github.com/tomruk/adsbrecon/sdr"
	"gopkg.in/abiosoft/ishell.v2"
)

func initCommands() {
	shell.AddCmd(&ishell.Cmd{
		Name: "choose-device",
		Help: "Choose the SDR device to listen on",
		Func: chooseDeviceCommand,
	})
}

func chooseDeviceCommand(c *ishell.Context) {
	device := chooseDevice(true)
	if device != nil {
		if sdr.Device != nil {
			sdr.Device.Stop()
		}

		err := device.Run()
		if err == nil {
			sdr.Device = device
		} else {
			fmt.Printf("Error occured while initializing the device: %s", err.Error())
		}
	}
}
