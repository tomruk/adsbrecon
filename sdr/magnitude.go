package sdr

import (
	"fmt"
)

/*
	Theese functions and calculations are
	implemented from: https://github.com/osmocom/rtl-sdr/blob/master/src/rtl_adsb.c
*/

const (
	BADSAMPLE      = 255
	OVERWRITE      = 254
	MESSAGEGO      = 253
	LONG_FRAME     = 112
	SHORT_FRAME    = 56
	PREAMBLE_LEN   = 16
	ALLOWED_ERRORS = 5

	VERBOSE = true
)

var (
	squares [256]uint16
	quality = 20
)

func init() {
	/*
		Precompute the squares
		Equivalent to abs(x-128) ^ 2
	*/
	var i uint16 = 0
	for ; i < 256; i++ {
		j := abs8(i)
		squares[i] = j * j
	}
}

func abs8(x uint16) uint16 {
	if x >= 127 {
		return x - 127
	} else {
		return 127 - x
	}
}

func magnitude(buf []uint8, len int) int {
	for i := 0; i < len; i += 2 {
		buf[i] = uint8(squares[buf[i]] + squares[buf[i]+1])
	}

	return len / 2
}

// Overwrites magnitude buffer with valid bits (BADSAMPLE on errors)
func manchester(buf []uint8, len int) {
	// a and b hold old values to verify the local manchester
	var (
		a            uint16 = 0
		b            uint16 = 0
		bit          uint16
		i            = 0
		j            int
		start        int
		errors       int
		maximumIndex = len - 1
	)

	for i < maximumIndex {
		// Find preamble
		for ; i < (len - PREAMBLE_LEN); i++ {
			if preamble(buf, i) == 0 {
				continue
			}

			a = uint16(buf[i])
			b = uint16(buf[i+1])
			for j = 0; j < PREAMBLE_LEN; j++ {
				buf[i+j] = MESSAGEGO
			}

			i += PREAMBLE_LEN
			break
		}

		start = i
		j = start
		errors = 0

		// Mark bits until encoding breaks
		for ; i < maximumIndex; i, j = i+2, j+1 {
			bit = singleManchester(a, b, uint16(buf[i]), uint16(buf[i+1]))
			a = uint16(buf[i])
			b = uint16(buf[i+1])
			if bit == BADSAMPLE {
				errors += 1
				if errors > ALLOWED_ERRORS {
					buf[j] = BADSAMPLE
					break
				} else {
					if a > b {
						bit = 1
					} else {
						bit = 0
					}

					// These don't have to match the bit
					a = 0
					b = 65535
				}
			}

			buf[i+1] = OVERWRITE
			buf[i] = buf[i+1]
			buf[j] = uint8(bit)
		}
	}
}

// Takes 4 consecutive real samples, return 0 or 1, BADSAMPLE on error
func singleManchester(a, b, c, d uint16) uint16 {
	var (
		bit  = c > d
		bitp = a > b
	)

	if quality == 0 {
		if bit {
			return 1
		} else {
			return 0
		}
	}

	if quality == 5 {
		if bit && bitp && b > c {
			return BADSAMPLE
		}

		if !bit && !bitp && b < c {
			return BADSAMPLE
		}

		if bit {
			return 1
		} else {
			return 0
		}
	}

	if quality == 10 {
		if bit && bitp && c > b {
			return 1
		}

		if bit && !bitp && d < b {
			return 1
		}

		if !bit && bitp && d > b {
			return 0
		}

		if !bit && !bitp && c < b {
			return 0
		}

		return BADSAMPLE
	}

	if bit && bitp && c > b && d < a {
		return 1
	}
	if bit && !bitp && c > a && d < b {
		return 1
	}
	if !bit && bitp && c < a && d > b {
		return 0
	}
	if !bit && !bitp && c < b && d > a {
		return 0
	}

	return BADSAMPLE
}

// Returns 0/1 for preamble at index i
func preamble(buf []uint8, i int) int {
	var (
		low  uint16 = 0
		high uint16 = 65535
	)

	for j := 0; j < PREAMBLE_LEN; j++ {
		switch j {
		case 0:
		case 2:
		case 7:
		case 9:
			high = uint16(buf[i+j])
		default:
			low = uint16(buf[i+j])
		}

		if high < low {
			return 0
		}
	}

	return 1
}

func dumpMessages(buf []uint8, len int) {
	var (
		i         = 0
		bufIndex  int
		index     int
		shift     uint
		frameLen  int
		adsbFrame [14]int
	)

	for ; i < len; i++ {
		if buf[i] > 1 {
			continue
		}

		frameLen = LONG_FRAME
		bufIndex = 0

		for index = 0; index < 14; index++ {
			adsbFrame[index] = 0
		}

		for ; i < len && buf[i] <= 1 && bufIndex < frameLen; i, bufIndex = i+1, bufIndex+1 {
			if buf[i] >= 1 {
				index = bufIndex / 8
				shift = (uint)(7 - (bufIndex % 8))

				shifted := 1 << shift
				adsbFrame[index] |= shifted
			}

			if bufIndex == 7 {
				if adsbFrame[0] == 0 {
					break
				}

				if (adsbFrame[0] & 0x80) >= 1 {
					frameLen = LONG_FRAME
				} else {
					frameLen = SHORT_FRAME
				}
			}
		}
		if bufIndex < (frameLen - 1) {
			continue
		}

		displayDumpedMessage(adsbFrame, frameLen)
	}
}

func displayDumpedMessage(adsbFrame [14]int, len int) {
	var (
		i  int
		df int
	)

	if len <= SHORT_FRAME {
		return
	}

	df = (adsbFrame[0] >> 3) & 0x1f
	if quality == 0 && !(df == 11 || df == 17 || df == 18 || df == 19) {
		return
	}

	fmt.Print("*")
	for i = 0; i < ((len + 7) / 8); i++ {
		fmt.Printf("%02x", adsbFrame[i])
	}

	fmt.Print(";\r\n")

	fmt.Printf("DF = %d CA = %d\n", df, adsbFrame[0]&0x07)
	fmt.Printf("ICAO Address = %06x\n", adsbFrame[1]<<16|adsbFrame[2]<<8|adsbFrame[3])

	if len <= SHORT_FRAME {
		return
	}

	fmt.Printf("PI = 0x%06x\n", adsbFrame[11]<<16|adsbFrame[12]<<8|adsbFrame[13])
	fmt.Printf("Type Code = %d S.Type/Ant. = %x\n", (adsbFrame[4]>>3)&0x1f, adsbFrame[4]&0x07)
	fmt.Print("--------------\n")
}
