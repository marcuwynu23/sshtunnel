package main

import (
    "log"
    "os"
    "path/filepath"
    "time"
    "github.com/marcuwynu23/sshtunnel/sshlib"
    "github.com/fsnotify/fsnotify"
)

func main() {
    // Setup logging to both file and console
    sshlib.SetupLogging()

    // Print banner at startup
    sshlib.PrintBanner()

    // Get the executable's directory and form the full config file path
    executablePath, err := os.Executable()
    if err != nil {
        log.Fatalf("Error retrieving executable path: %v", err)
    }
    executableDir := filepath.Dir(executablePath)
    configFilePath := filepath.Join(executableDir, "sshtunnel.yml")

    // Load the initial config
    config, err := sshlib.LoadConfig(configFilePath)
    if err != nil {
        log.Fatalf("Error loading config file: %v", err)
    }

    // Start the SSH tunneling
    log.Println("Configuration loaded successfully.")
    startSSHTunneling(config)

    // Monitor config.yml for changes
    watchConfigFile(configFilePath)
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
