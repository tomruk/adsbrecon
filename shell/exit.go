package shell

import (
	"fmt"
	"os"
	"time"

	"github.com/tomruk/adsbrecon/sdr"
)

var (
	exitChannel = make(chan bool)
)

func exit(code int) {
	fmt.Println()

	go killall()
	<-exitChannel

	fmt.Println("Exiting")
	os.Exit(code)
}

// Kills all running threads
// TODO: Implement this function
func killall() {
	go timeout()

	if sdr.Device != nil {
		sdr.Device.Stop()
	}

	exitChannel <- true
}

func timeout() {
	time.Sleep(time.Second * 2)
	exitChannel <- true
}
