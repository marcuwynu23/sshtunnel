package main

import (
	"log"
	"ssh-tunnel/sshlib"
)

func main() {
	// Setup logging to both file and console
	sshlib.SetupLogging()

	// Print banner at startup
	sshlib.PrintBanner()

	// Load config from YAML
	configFilePath := "config.yml"
	config, err := sshlib.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	log.Println("Configuration loaded successfully.")
	log.Println("Starting SSH tunneling...")

	// Keep the SSH connection alive and auto-reconnect on failure
	sshlib.MaintainSSHConnection(config)
}
