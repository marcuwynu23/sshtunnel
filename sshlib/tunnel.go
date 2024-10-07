package sshlib

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type Config struct {
	SSHConfig SSHConfig `yaml:"ssh_config"`
}

type SSHConfig struct {
	Host       string   `yaml:"host"`
	Port       int      `yaml:"port"`
	User       string   `yaml:"user"`
	PrivateKey string   `yaml:"private_key"`
	Tunnels    []Tunnel `yaml:"tunnels"`
}

type Tunnel struct {
	LocalIP    string `yaml:"local_ip"`   // New field to specify local IP
	LocalPort  int    `yaml:"local_port"`
	RemotePort int    `yaml:"remote_port"`
}

// Function to load YAML config
func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// Function to dial SSH connection
func sshDial(config *SSHConfig) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	sshAddr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	client, err := ssh.Dial("tcp", sshAddr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %v", err)
	}
	return client, nil
}

// Function to start the reverse SSH tunnel
func startTunnel(client *ssh.Client, config *SSHConfig, tunnel Tunnel) error {
    // Listen on the remote server's localhost for the tunnel
    listener, err := client.Listen("tcp", fmt.Sprintf("localhost:%d", tunnel.RemotePort))
    if err != nil {
        return fmt.Errorf("Failed to set up remote listener: %v", err)
    }
    defer listener.Close()

    // Log the host address and tunnel details
    log.Printf("Tunneling from remote address %s:%d to local address %s:%d", config.Host, tunnel.RemotePort, tunnel.LocalIP, tunnel.LocalPort)

    for {
        conn, err := listener.Accept()
        if err != nil {
            return fmt.Errorf("Listener accept failed: %v", err)
        }

        go handleTunnel(conn, tunnel.LocalIP, tunnel.LocalPort)
    }
}


// Function to handle individual tunnel connections
func handleTunnel(conn net.Conn, localIP string, localPort int) {
    defer conn.Close()

    // Use the LocalIP and LocalPort to dial the local service
    localConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", localIP, localPort))
    if err != nil {
        log.Printf("Failed to connect to local service at %s:%d: %v", localIP, localPort, err)
        return
    }
    defer localConn.Close()

    log.Printf("Connection established to local service at %s:%d", localIP, localPort)

    // Forward traffic between local service and remote connection
    go func() { _, _ = io.Copy(localConn, conn) }()
    _, _ = io.Copy(conn, localConn)
}


// Function to continuously try reconnecting if the SSH connection fails
func MaintainSSHConnection(config *Config) {
	for {
		client, err := sshDial(&config.SSHConfig)
		if err != nil {
			log.Printf("Failed to establish SSH connection: %v. Retrying in 10 seconds...", err)
			time.Sleep(10 * time.Second) // Wait before retrying
			continue
		}

		// Start tunnels once connected
		for _, tunnel := range config.SSHConfig.Tunnels {
			go func(tunnel Tunnel) {
				for {
					err := startTunnel(client, &config.SSHConfig, tunnel)
					if err != nil {
						log.Printf("Tunnel failed: %v. Reconnecting...", err)
						break // Exit the loop and retry the connection
					}
				}
			}(tunnel)
		}

		// Keep the SSH connection alive
		client.Wait()
		log.Printf("SSH connection lost. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
}

// Set up logging
func SetupLogging() {
	// Open a log file in append mode
	logFile, err := os.OpenFile("ssh_tunneling.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create a multi-writer to write logs both to file and to console (optional)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set output of the standard logger to both file and console
	log.SetOutput(multiWriter)
}

// Print the startup banner
func PrintBanner() {
	fmt.Println("=======================================")
	fmt.Println("   SSH Tunneling Service")
	fmt.Println("   Securely tunnel your services")
	fmt.Println("=======================================")
}
