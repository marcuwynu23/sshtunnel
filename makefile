# Project name
BINARY_NAME=sshtunnel

# Build output directory
BUILD_DIR=build

# Use VERSION if provided (CI), otherwise latest tag, fallback to dev
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo dev)

dev:
	go build -o sshtunnel.exe -buildvcs=false
	xcopy /Y /C sshtunnel.exe D:\Executables\sshtunnel

# Go architectures and operating systems
ARCH_LIST=amd64 386 arm64 arm
OS_LIST=linux windows darwin

# YAML configuration file
CONFIG_FILE=sshtunnel.yml

# Default target to build for all OS and architecture combinations
.PHONY: dev ci all clean clean-build linux_amd64 linux_386 windows_amd64 windows_386 macos_amd64 linux_arm64 linux_arm done build

ci:
	go fmt ./...
	go vet ./...
	go test ./...
	go build ./...

all: clean linux_amd64 windows_amd64 macos_amd64 linux_386 windows_386 linux_arm64 linux_arm clean-build

# Clean up the build directory
clean:
	rm -rf $(BUILD_DIR)
clean-build:
	find $(BUILD_DIR) -type f ! -name "*.tar.gz" -exec rm -f {} \;
# Linux 64-bit build (amd64)
linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_amd64_$(VERSION).tar $(BINARY_NAME)_linux_amd64 $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64_$(VERSION).tar
	
# Linux 32-bit build (386)
linux_386:
	GOOS=linux GOARCH=386 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_386 -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_386_$(VERSION).tar $(BINARY_NAME)_linux_386 $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_linux_386_$(VERSION).tar


# Windows 64-bit build (amd64)
windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_windows_amd64_$(VERSION).tar $(BINARY_NAME)_windows_amd64.exe $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64_$(VERSION).tar


# Windows 32-bit build (386)
windows_386:
	GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/$(BINARY_NAME)_windows_386.exe -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_windows_386_$(VERSION).tar $(BINARY_NAME)_windows_386.exe $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_windows_386_$(VERSION).tar
	

# MacOS 64-bit build (amd64 only, no 32-bit support)
macos_amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_macos_amd64 -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_macos_amd64_$(VERSION).tar $(BINARY_NAME)_macos_amd64 $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_macos_amd64_$(VERSION).tar
	

# ARM 64-bit build for Linux (arm64)
linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_arm64_$(VERSION).tar $(BINARY_NAME)_linux_arm64 $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64_$(VERSION).tar
	

# ARM 32-bit build for Linux (arm)
linux_arm:
	GOOS=linux GOARCH=arm go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm -buildvcs=false main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_arm_$(VERSION).tar $(BINARY_NAME)_linux_arm $(CONFIG_FILE)
	gzip $(BUILD_DIR)/$(BINARY_NAME)_linux_arm_$(VERSION).tar
	
# Message when builds are complete
done:
	@echo "Builds and tar archives for multiple architectures and OS complete! Version: $(VERSION)"

# Full build process
build: clean linux_amd64 windows_amd64 macos_amd64 linux_386 windows_386 linux_arm64 linux_arm done
