# **SSHTunnel**

**SSHTunnel** is a cross-platform command-line utility designed for **reverse SSH tunneling** using a configuration file (`sshtunnel.yml`). This tool simplifies the process of creating SSH tunnels by using a predefined configuration, making it easy to manage multiple tunnels without specifying all the parameters in the command line.

## **Key Features**

- **Reverse SSH Tunneling**: Facilitates SSH tunneling from a remote machine back to the local machine, enabling access to services behind NAT or firewalls.
- **YAML Configuration**: Configure multiple tunnels and SSH options through a simple YAML file.
- **Cross-Platform**: Works on Linux, Windows, and macOS.
- **Multiple Tunnels**: Supports defining and running multiple tunnels at once through the configuration file.
- **Secure Authentication**: Uses private keys for SSH authentication.

## **What is Reverse SSH Tunneling?**

**Reverse SSH tunneling** allows a machine that is behind a firewall or without a public IP to expose a port on a remote server that has a public IP. This technique is useful for:

- Accessing devices or servers behind firewalls or NAT.
- Exposing local services (e.g., web servers, databases) securely to the internet.
- Remotely managing devices or systems without direct public access.

## **How It Works**

1. **SSHTunnel** reads from a configuration file (`sshtunnel.yml`) to establish an SSH connection to a remote server.
2. The remote server sets up reverse tunnels as defined in the configuration file.
3. Any connection to the specified remote ports on the server will be tunneled to the corresponding local ports on the remote machine.

### **Example Configuration**

```yaml
ssh_config:
  host: "<hostip>"
  port: 22
  user: "<remote_username>"
  private_key: "C:\Users\<local_username>\.ssh\id_rsa"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 5000
      remote_port: 7000
    - local_ip: "0.0.0.0"
      local_port: 5200
      remote_port: 7200
```

- **host**: The IP address or hostname of the remote server.
- **port**: The SSH port (usually `22`).
- **user**: The SSH username for the remote server.
- **private_key**: The path to the private SSH key used for authentication.
- **tunnels**: An array of tunnels where each tunnel forwards traffic from a remote port on the server to a local port on the remote machine.

### **Example Use Case**

Let’s say you want to access a service running on port `5000` on your local machine, which is behind a NAT firewall, and make it available on port `7000` of a public server. Additionally, you want to forward another local port (`5200`) to the public server's port `7200`. You would define this configuration in `sshtunnel.yml`:

```yaml
ssh_config:
  host: "123.45.67.89"
  port: 22
  user: "remote_user"
  private_key: "/home/user/.ssh/id_rsa"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 5000
      remote_port: 7000
    - local_ip: "0.0.0.0"
      local_port: 5200
      remote_port: 7200
```

When **SSHTunnel** is executed, it will set up two reverse tunnels:

1. Remote port `7000` on the public server will forward to local port `5000` on the local machine.
2. Remote port `7200` on the public server will forward to local port `5200` on the local machine.

## **Installation**

### **Download Prebuilt Binaries**

You can download prebuilt binaries for your platform from the [releases page](https://github.com/marcuwynu23/sshtunnel/releases).

Available binaries:

- `sshtunnel_linux_amd64`
- `sshtunnel_windows_amd64.exe`
- `sshtunnel_macos_amd64`

### **Building from Source**

To build SSHTunnel from source, you need to have Go installed on your system. Follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/sshtunnel.git
   cd sshtunnel
   ```

2. Build the binary:

   ```bash
   go build -o sshtunnel main.go
   ```

3. The binary `sshtunnel` will be generated in your current directory.

## **Usage**

### **1. Create the Configuration File**

Create a configuration file called `sshtunnel.yml` with your SSH connection details and the tunnels you want to establish. For example:

```yaml
ssh_config:
  host: "123.45.67.89"
  port: 22
  user: "remote_user"
  private_key: "/home/user/.ssh/id_rsa"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 5000
      remote_port: 7000
    - local_ip: "0.0.0.0"
      local_port: 5200
      remote_port: 7200
```

### **2. Run the SSHTunnel Command**

Once the configuration file is set up, you can start the SSH tunnel by simply running:

```bash
sshtunnel
```

This will read the `sshtunnel.yml` file from the current directory and establish the reverse SSH tunnels as defined in the configuration.

### **3. Configuration Options**

- **host**: The IP address or domain of the SSH server.
- **port**: The SSH port (default is `22`).
- **user**: SSH username on the remote server.
- **private_key**: Path to your private SSH key for authentication.
- **tunnels**: Define an array of tunnels. Each tunnel requires:
  - `local_ip`: The local IP (usually `0.0.0.0` to bind to all interfaces).
  - `local_port`: The port on the local machine to forward.
  - `remote_port`: The port on the remote server to expose the forwarded service.

### **Running in the Background**

To run SSHTunnel as a background service, you can use your system's process control methods (e.g., `nohup`, `systemd`, or `Windows Services`).

### **Examples**

1. **Expose a local web server to the public:**

   In the `sshtunnel.yml`:

   ```yaml
   ssh_config:
     host: "example.com"
     port: 22
     user: "remote_user"
     private_key: "/home/user/.ssh/id_rsa"
     tunnels:
       - local_ip: "0.0.0.0"
         local_port: 8080
         remote_port: 9090
   ```

   Running `sshtunnel` will forward local port `8080` to port `9090` on `example.com`.

2. **Expose SSH from behind a NAT:**

   To expose your SSH port behind a NAT to a remote server:

   ```yaml
   ssh_config:
     host: "123.45.67.89"
     port: 22
     user: "remote_user"
     private_key: "~/.ssh/id_rsa"
     tunnels:
       - local_ip: "0.0.0.0"
         local_port: 22
         remote_port: 2222
   ```

   Running `sshtunnel` will expose your local SSH on port `22` to the remote server on port `2222`.

## **System Requirements**

- **Go 1.16+** for building from source.
- **SSH** client installed on your system.
- Binaries available for the following platforms:
  - Linux (64-bit and 32-bit)
  - Windows (64-bit and 32-bit)
  - macOS (64-bit)

## **Contributing**

We welcome contributions! If you find a bug or have a feature request, feel free to open an issue or submit a pull request.

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Open a pull request.

## **License**

SSHTunnel is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
