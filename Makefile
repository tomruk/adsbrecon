GO = go
PREFIX ?= /usr/local/bin

ARM_CC ?= arm-linux-gnueabihf-gcc
ARM_CXX ?= arm-linux-gnueabihf-g++

ARM64_CC ?= aarch64-linux-gnu-gcc
ARM64_CXX ?= aarch64-linux-gnu-g++

NDK_ROOT ?= /opt/android-sdk/ndk-bundle
NDK_CC ?= $(NDK_ROOT)/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang
NDK_CXX ?= $(NDK_ROOT)/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android28-clang++

MINGW_64_CC ?= x86_64-w64-mingw32-gcc
MINGW_86_CC ?= i686-w64-mingw32-gcc
MINGW_86_CXX ?= i686-w64-mingw32-g++
MINGW_64_CXX ?= x86_64-w64-mingw32-g++

all: build-native

.PHONY: install
install:
	mv adsbrecon $(PREFIX)

.PHONY: build-all
build-all: build-linux-x64 build-linux-x86 build-linux-arm64 build-linux-armv5 build-linux-armv6 build-linux-armv7 build-android build-windows-x64 build-windows-x86

.PHONY: build-native
build-native:
	CGO_ENABLED=1 $(GO) build -o adsbrecon

.PHONY: build-linux-x64
build-linux-x64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 $(GO) build -o adsbrecon-linux-x64

.PHONY: build-linux-x86
build-linux-x86:
	GOOS=linux GOARCH=386 CGO_ENABLED=1 $(GO) build -o adsbrecon-linux-x86

.PHONY: build-linux-arm64
build-linux-arm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=$(ARM64_CC) CXX=$(ARM64_CXX) $(GO) build -o adsbrecon-linux-arm64

.PHONY: build-linux-armv5
build-linux-armv5:
	GOOS=linux GOARCH=arm GOARM=5 CGO_ENABLED=1 CC=$(ARM_CC) CXX=$(ARM_CXX) $(GO) build -o adsbrecon-linux-armv5

.PHONY: build-linux-armv6
build-linux-armv6:
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=$(ARM_CC) CXX=$(ARM_CXX) $(GO) build -o adsbrecon-linux-armv6

.PHONY: build-linux-armv7
build-linux-armv7:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=$(ARM_CC) CXX=$(ARM_CXX) $(GO) build -o adsbrecon-linux-armv7

.PHONY: build-android
build-android:
	GOOS=android GOARCH=arm64 GOARM=7 CGO_ENABLED=1 CC=$(NDK_CC) CXX=$(NDK_CXX) $(GO) build -o adsbrecon-android

.PHONY: build-windows-x64
build-windows-x64:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=$(MINGW_64_CC) CXX=$(MINGW_64_CXX) $(GO) build -o adsbrecon-windows-x64.exe

.PHONY: build-windows-x86
build-windows-x86:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=$(MINGW_86_CC) CXX=$(MINGW_86_CXX) $(GO) build -o adsbrecon-windows-x86.exe

.PHONY: clean
clean:
	rm -f adsbrecon adsbrecon.exe adsbrecon-linux-x64 adsbrecon-linux-x86 adsbrecon-linux-arm64 adsbrecon-linux-armv7 adsbrecon-linux-armv6 adsbrecon-linux-armv5 adsbrecon-windows-x64.exe adsbrecon-windows-x86.exe
