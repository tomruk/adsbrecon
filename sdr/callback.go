package sdr

import (
	"sync"
)

const (
	DATA_LEN = 16 * 32 * 512
)

var (
	dataMutex   = &sync.Mutex{}
	dataChannel = make(chan bool)

	data [DATA_LEN]uint8
)

func init() {
	go demodulationThread()
}

func callback(buf []byte) {
	dataMutex.Lock()
	dataChannel <- true

	copy(data[:], buf)

	dataMutex.Unlock()
}

func demodulationThread() {
	for {
		<-dataChannel
		dataMutex.Lock()

		len := magnitude(data[:], DATA_LEN)
		manchester(data[:], len)
		dumpMessages(data[:], len)

		dataMutex.Unlock()
	}
}
