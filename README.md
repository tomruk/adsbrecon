# ADS-B Recon

This is an attempt to make an ADS-B receiver software using Go.  
Inspired from dump1090 and rtl_adsb.  

## Warning

This project is neither complete nor i have time to write it comprehensively. **Don't use it**.

## Supported hardware

**RTL-SDR** -> Yes  
**LimeSDR** -> No  
**HackRF** -> No  

## Usage

Install [librtlsdr](https://github.com/steve-m/librtlsdr), then:

```shell
make
sudo make install
```