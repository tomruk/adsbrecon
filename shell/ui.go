package shell

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/tomruk/adsbrecon/sdr"
	"gopkg.in/abiosoft/ishell.v2"
)

var (
	shell *ishell.Shell

	underline = color.New(color.Underline)
	blue      = color.New(color.FgBlue)
)

func init() {
	shell = ishell.New()
	shell.SetOut(color.Output)
}

func Run() {
	splash()

	// CTRL+C handler
	shell.Interrupt(func(c *ishell.Context, count int, _ string) {
		shellInterrupt()
	})

	// CTRL+D handler
	shell.EOF(func(c *ishell.Context) {
		shellInterrupt()
	})

	device := deviceSelection()
	if device != nil {
		err := device.Run()
		if err == nil {
			sdr.Device = device
		} else {
			fmt.Printf("Error occured when intializing the device: %s", err.Error())
		}
	}

	fmt.Println()
	initCommands()
	updatePrompt("")
	shell.Run()
}

func updatePrompt(prompt string) {
	adsbreconText := underline.Sprint("adsbrecon")

	if prompt == "" {
		shell.SetPrompt(fmt.Sprintf("%s » ", adsbreconText))
	} else {
		shell.SetPrompt(fmt.Sprintf("%s » ", prompt))
	}
}

func shellInterrupt() {
	if confirmation("Are you sure you want to exit?", false) {
		exit(0)
	} else {
		fmt.Println("Let's move on")
	}
}

func confirmation(prompt string, yesDefault bool) bool {
	shell.ShowPrompt(false)
	defer shell.ShowPrompt(true)

	if yesDefault {
		prompt = prompt + " Y/n "
	} else {
		prompt = prompt + " y/N "
	}

	shell.Print(prompt)

	confirmation := strings.TrimSpace(strings.ToLower(shell.ReadLine()))

	if (confirmation == "y") || (yesDefault && confirmation == "") {
		return true
	} else {
		return false
	}
}

func deviceSelection() (device sdr.SDRDevice) {
	device = nil

	if *arguments.deviceType == "" {
		device = chooseDevice(false)
		return
	} else {
		found := false

		switch *arguments.deviceType {
		case "rtlsdr":
			rtlsdrDevices := sdr.GetRTLSDRDevices(0)

			for _, d := range rtlsdrDevices {
				if d.GetIndex() == *arguments.deviceIndex {
					found = true
					device = d
					return
				}
			}
		}

		if !found {
			fmt.Printf("Device with index #%d not found\n", *arguments.deviceIndex)
		}
	}

	return
}

func splash() {
	radar := `
		  
          ,-.
         / \  '.  __..-,-
        :   \ --''_..-'.'
        |    . .-' '. '.
        :     .     .'.'
         \     '.  /  ..
          \      '.   ' .
           ',       '.   \
          ,|,'.        '-.\
         '.||  ''-...__..-'
          |  |
          |__|
          /||\
         //||\\
        // || \\
     __//__||__\\__
    '--------------'
			

`

	blue.Print(radar + "\r")
}
