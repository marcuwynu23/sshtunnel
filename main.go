package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/marcuwynu23/sshtunnel/sshlib"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Get the executable's directory and form the default config file path.
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error retrieving executable path: %v", err)
	}
	executableDir := filepath.Dir(executablePath)
	defaultConfigFilePath := filepath.Join(executableDir, "sshtunnel.yml")

	// CLI options.
	configFlag := flag.String("config", "", "Path to config file (default: executable_dir/sshtunnel.yml)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [--config <config-file>]\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	configFilePath := defaultConfigFilePath
	if *configFlag != "" {
		configFilePath = *configFlag
	}

	absConfigFilePath, err := filepath.Abs(configFilePath)
	if err != nil {
		log.Fatalf("Error resolving config file path: %v", err)
	}

	// Setup logging to both file and console, with log path near config.
	sshlib.SetupLogging(filepath.Dir(absConfigFilePath))

	// Print banner at startup
	sshlib.PrintBanner()

	// Load the initial config
	config, err := sshlib.LoadConfig(absConfigFilePath)
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}

	// Start the SSH tunneling
	log.Println("Configuration loaded successfully.")
	startSSHTunneling(config)

	// Monitor config.yml for changes
	watchConfigFile(absConfigFilePath)
}

func startSSHTunneling(config *sshlib.Config) {
	log.Println("Starting SSH tunneling...")

	// Function to maintain SSH connection
	go sshlib.MaintainSSHConnection(config)
}

func watchConfigFile(filePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start watching the file
	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching config file for changes...")

	// Listen for file system events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("Config file changed. Reloading...")

				// Add a small delay to ensure the file has finished writing
				time.Sleep(500 * time.Millisecond)

				// Reload the configuration
				config, err := sshlib.LoadConfig(filePath)
				if err != nil {
					log.Printf("Error reloading config file: %v", err)
				} else {
					log.Println("Configuration reloaded successfully.")

					// Restart SSH tunneling with the new configuration
					startSSHTunneling(config)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error watching config file:", err)
		}
	}
}
