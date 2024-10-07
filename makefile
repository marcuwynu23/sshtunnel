# Project name
BINARY_NAME=sshtunnel

# Build output directory
BUILD_DIR=build

# Fetch the latest Git tag, stripping any extra commit info
GIT_TAG=$(shell git describe --tags --abbrev=0)

# Go architectures and operating systems
ARCH_LIST=amd64 386 arm64 arm
OS_LIST=linux windows darwin

# YAML configuration file
CONFIG_FILE=sshtunnel.yml

# Default target to build for all OS and architecture combinations
all: clean linux_amd64 windows_amd64 macos_amd64 linux_386 windows_386 linux_arm64 linux_arm

# Clean up the build directory
clean:
	rm -rf $(BUILD_DIR)

# Linux 64-bit build (amd64)
linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_amd64_$(GIT_TAG).tar $(BINARY_NAME)_linux_amd64 $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# Linux 32-bit build (386)
linux_386:
	GOOS=linux GOARCH=386 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_386 main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_386_$(GIT_TAG).tar $(BINARY_NAME)_linux_386 $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_linux_386
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# Windows 64-bit build (amd64)
windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_windows_amd64_$(GIT_TAG).tar $(BINARY_NAME)_windows_amd64.exe $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# Windows 32-bit build (386)
windows_386:
	GOOS=windows GOARCH=386 go build -o $(BUILD_DIR)/$(BINARY_NAME)_windows_386.exe main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_windows_386_$(GIT_TAG).tar $(BINARY_NAME)_windows_386.exe $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_windows_386.exe
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# MacOS 64-bit build (amd64 only, no 32-bit support)
macos_amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_macos_amd64 main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_macos_amd64_$(GIT_TAG).tar $(BINARY_NAME)_macos_amd64 $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_macos_amd64
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# ARM 64-bit build for Linux (arm64)
linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_arm64_$(GIT_TAG).tar $(BINARY_NAME)_linux_arm64 $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# ARM 32-bit build for Linux (arm)
linux_arm:
	GOOS=linux GOARCH=arm go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm main.go
	cp $(CONFIG_FILE) $(BUILD_DIR)
	cd $(BUILD_DIR) && tar -cvf $(BINARY_NAME)_linux_arm_$(GIT_TAG).tar $(BINARY_NAME)_linux_arm $(CONFIG_FILE)
	rm $(BUILD_DIR)/$(BINARY_NAME)_linux_arm
	rm $(BUILD_DIR)/$(CONFIG_FILE)

# Message when builds are complete
done:
	@echo "Builds and tar archives for multiple architectures and OS complete! Version: $(GIT_TAG)"

# Full build process
build: clean linux_amd64 windows_amd64 macos_amd64 linux_386 windows_386 linux_arm64 linux_arm done
