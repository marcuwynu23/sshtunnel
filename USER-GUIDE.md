# SSHTunnel User Guide

➡️ **[Back to README →](README.md)**

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Concepts](#concepts)
- [CI/CD Integration](#cicd-integration)
- [Workflows](#workflows)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Installation

### Prerequisites

- **Prebuilt binary**: None — the binary is self-contained.
- **Building from source**: [Go 1.23+](https://go.dev/dl/)

### Download a Release

1. Visit the [releases page](https://github.com/marcuwynu23/sshtunnel/releases).
2. Download the archive for your platform:

   | Platform | Architecture | File pattern |
   |---|---|---|
   | Linux | amd64 | `sshtunnel_linux_amd64_<version>.tar.gz` |
   | Linux | 386 | `sshtunnel_linux_386_<version>.tar.gz` |
   | Linux | arm64 | `sshtunnel_linux_arm64_<version>.tar.gz` |
   | Linux | arm | `sshtunnel_linux_arm_<version>.tar.gz` |
   | Windows | amd64 | `sshtunnel_windows_amd64_<version>.tar.gz` |
   | Windows | 386 | `sshtunnel_windows_386_<version>.tar.gz` |
   | macOS | amd64 | `sshtunnel_macos_amd64_<version>.tar.gz` |

3. Extract and run:

   ```bash
   tar xzf sshtunnel_linux_amd64_v1.2.3.tar.gz
   cd sshtunnel_linux_amd64_v1.2.3
   ./sshtunnel
   ```

### Build from Source

```bash
git clone https://github.com/marcuwynu23/sshtunnel.git
cd sshtunnel
go build -o sshtunnel main.go
./sshtunnel --help
```

### Verify

```bash
./sshtunnel --help
```

Expected output:

```
Usage: sshtunnel [--config <config-file>]
Options:
  --config string
        Path to config file (default: executable_dir/sshtunnel.yml)
```

---

## Quick Start

### 1. Create a config file

Save this as `sshtunnel.yml` in the same directory as the binary:

```yaml
ssh_config:
  host: "203.0.113.1"
  port: 22
  user: "deploy"
  private_key: "/home/deploy/.ssh/id_ed25519"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 8080
      remote_ip: "0.0.0.0"
      remote_port: 9090
```

### 2. Run the tool

```bash
./sshtunnel
```

### 3. Verify the tunnel

From the remote server, connect to the exposed port:

```bash
curl http://localhost:9090
```

Traffic reaches the service running on `localhost:8080` of the machine that initiated the tunnel.

---

## Command Reference

### `sshtunnel`

Start the SSH tunneling service.

```bash
sshtunnel [--config <path>]
```

| Flag | Default | Description |
|---|---|---|
| `--config` | `<executable_dir>/sshtunnel.yml` | Path to the YAML configuration file |
| `--help` | — | Print usage and exit |

#### Examples by Use Case

**Use the default config (next to the binary):**

```bash
sshtunnel
```

The tool looks for `sshtunnel.yml` in the same directory as the executable and creates `ssh_tunneling.log` there.

**Use a custom config file:**

```bash
sshtunnel --config /etc/sshtunnel/production.yml
```

Logs are written to `/etc/sshtunnel/ssh_tunneling.log`.

**Show usage:**

```bash
sshtunnel --help
```

---

## Configuration

### Full YAML Schema

```yaml
ssh_config:
  # Required — remote SSH server
  host: "203.0.113.1"

  # Optional — default: 22
  port: 22

  # Required — SSH login user
  user: "deploy"

  # Required — path to the private key file
  private_key: "/home/deploy/.ssh/id_ed25519"

  # Required — one or more tunnel definitions
  tunnels:
    - local_ip: "0.0.0.0"      # Local bind IP
      local_port: 8080           # Local port to forward traffic to
      remote_ip: "0.0.0.0"      # Remote bind IP on the SSH server
      remote_port: 9090          # Remote port to expose

    - local_ip: "127.0.0.1"
      local_port: 3000
      remote_ip: "0.0.0.0"
      remote_port: 3000
```

### Field Reference

| Field | Type | Required | Default | Description |
|---|---|---|---|---|
| `ssh_config.host` | string | yes | — | Remote SSH server hostname or IP address |
| `ssh_config.port` | int | no | `22` | TCP port for the SSH connection |
| `ssh_config.user` | string | yes | — | SSH authentication username |
| `ssh_config.private_key` | string | yes | — | Absolute or relative path to the private key file (must be in OpenSSH format) |
| `ssh_config.tunnels` | array | yes | `[]` | List of tunnel mappings |
| `tunnels[].local_ip` | string | no | `"0.0.0.0"` | Local IP address to bind for incoming forwarded connections |
| `tunnels[].local_port` | int | yes | — | Local port number to receive forwarded traffic |
| `tunnels[].remote_ip` | string | no | `"0.0.0.0"` | Remote IP address on which the SSH server listens |
| `tunnels[].remote_port` | int | yes | — | Remote port number exposed on the SSH server |

### Config Precedence

1. `--config` CLI flag (highest priority)
2. `<executable_dir>/sshtunnel.yml` (default)

The log file `ssh_tunneling.log` is always created in the same directory as the resolved config file.

### Config Hot-Reload

While SSHTunnel is running, edit and save the config file. `fsnotify` detects the `Write` event, re-reads the YAML, and restarts all tunnels with zero downtime. Invalid YAML is reported in the log; the previous configuration continues to run.

---

## Concepts

### Reverse SSH Tunneling

A normal (local) SSH tunnel forwards traffic from a local port to a remote destination. A **reverse** tunnel does the opposite: the SSH server listens on a remote port and forwards incoming connections back through the SSH session to a local service.

```
                    SSH connection (outbound)
  Local machine  ──────────────────────────────>  Remote server
  (behind NAT)                                     (public IP)
       │                                                  │
       │  local:8080                            remote:9090 │
       │  (your app)                          (public port) │
       └────────────  reverse tunnel  ←─────────────────────┘
                 traffic flows back through SSH
```

Use cases:

- Expose a webhook receiver behind a firewall
- Access a database on a machine without a public IP
- Manage IoT devices on private networks

### SSH Authentication

SSHTunnel authenticates via **private key only**. The key must be in OpenSSH format (generated by `ssh-keygen`). Password authentication and SSH agent forwarding are not supported.

The corresponding public key must be installed in the remote user's `~/.ssh/authorized_keys`.

### Connection Lifecycle

1. **Dial** — Connect to `host:port` as `user`, authenticating with `private_key`.
2. **Listen** — For each tunnel, instruct the SSH server to listen on `remote_ip:remote_port`.
3. **Forward** — When a connection arrives at the remote listener, SSHTunnel dials `local_ip:local_port` and copies data bidirectionally.
4. **Retry** — If the SSH connection drops, wait 10 seconds and return to step 1.
5. **Hot-reload** — If the config file changes, tear down all listeners and return to step 1 with the new config.

---

## CI/CD Integration

### GitHub Actions

#### CI — Run on push / PR

```yaml
name: CI
on:
  push:
    branches: ["main"]
  pull_request:
    branches: ["main"]

permissions:
  contents: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: make ci
```

#### Release — Build and publish on tag

```yaml
name: Release
on:
  push:
    tags: ["v*"]

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      - run: make build VERSION=${GITHUB_REF_NAME}
      - uses: softprops/action-gh-release@v2
        with:
          files: build/*.tar.gz
```

### Running SSHTunnel in CI

If your CI runner needs to establish a tunnel (e.g., for integration tests), run it in the background:

```bash
# Start tunnel in background
./sshtunnel --config ./ci-sshtunnel.yml &

# Wait for tunnel to establish
sleep 3

# Run your tests
go test ./...

# Clean up
kill %1
```

---

## Workflows

### Expose a Local Web Server

```yaml
ssh_config:
  host: "vps.example.com"
  port: 22
  user: "app"
  private_key: "/home/app/.ssh/id_ed25519"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 8080
      remote_ip: "0.0.0.0"
      remote_port: 9090
```

Run `sshtunnel`. Users hitting `http://vps.example.com:9090` see your local web app.

### Remote SSH Access Behind NAT

```yaml
ssh_config:
  host: "vps.example.com"
  port: 22
  user: "admin"
  private_key: "/home/admin/.ssh/id_ed25519"
  tunnels:
    - local_ip: "127.0.0.1"
      local_port: 22
      remote_ip: "0.0.0.0"
      remote_port: 2222
```

From the VPS, connect back to the NATed machine:

```bash
ssh -p 2222 admin@localhost
```

### Multi-Tunnel Setup

```yaml
ssh_config:
  host: "vps.example.com"
  port: 22
  user: "deploy"
  private_key: "/home/deploy/.ssh/id_ed25519"
  tunnels:
    - local_ip: "127.0.0.1"
      local_port: 3000
      remote_port: 3000
    - local_ip: "127.0.0.1"
      local_port: 5432
      remote_port: 15432
    - local_ip: "0.0.0.0"
      local_port: 9090
      remote_port: 9090
```

Exposes three services: an app, a database, and a metrics endpoint.

### Running as a System Service (Linux)

Create a `systemd` unit:

```
[Unit]
Description=SSHTunnel
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/sshtunnel --config /etc/sshtunnel.yml
Restart=always
RestartSec=10
User=sshtunnel

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable sshtunnel
sudo systemctl start sshtunnel
```

### Running as a Background Process (Windows)

```cmd
start /B sshtunnel.exe --config C:\sshtunnel\sshtunnel.yml
```

Or configure as a Windows Scheduled Task / NSSM service.

---

## Troubleshooting

### Connection refused / timeout

| Cause | Fix |
|---|---|
| Remote SSH server is down or firewalled | Verify `ssh user@host -p port` works manually |
| Wrong host / port in config | Check `ssh_config.host` and `ssh_config.port` |
| Private key path is incorrect | Use an absolute path. Verify the file exists and has `0600` permissions |

### Permission denied (publickey)

| Cause | Fix |
|---|---|
| Public key not in `authorized_keys` | Run `ssh-copy-id user@host` to install it |
| Wrong key file | Confirm `private_key` points to the correct private key |
| Key format not supported | Generate a key with `ssh-keygen -t ed25519` (OpenSSH format) |

### Tunnel not accepting connections

| Cause | Fix |
|---|---|
| `remote_ip` is not `0.0.0.0` or a server-owned IP | Set `remote_ip: "0.0.0.0"` to bind on all interfaces |
| Remote port is already in use | Change `remote_port` or free the port on the server |
| `GatewayPorts no` on SSH server | Set `GatewayPorts yes` in `/etc/ssh/sshd_config` on the remote server |

### Config hot-reload not working

| Cause | Fix |
|---|---|
| File saved to a new inode (some editors) | Use `:w` / save-in-place instead of "Save As" |
| Permission denied watching config | Ensure the binary can read the config file |
| YAML syntax error | Check the log for `Error reloading config file` messages |

### Log file not created

| Cause | Fix |
|---|---|
| Directory not writable | Ensure the directory containing the config file is writable |
| Custom `--config` path | The log file is created in the **same directory** as the config file argument |

---

## FAQ

**Q: Do I need an SSH client installed on the machine running SSHTunnel?**

A: No. SSHTunnel uses the `golang.org/x/crypto/ssh` library — it is a pure Go SSH client compiled into the binary.

**Q: Can I use password authentication?**

A: No. Only private key authentication is supported (for security reasons).

**Q: Does SSHTunnel support SSH agent forwarding?**

A: Not currently. The private key must be readable by the binary as a file.

**Q: How do I add or remove tunnels without restarting?**

A: Edit the config file and save it. SSHTunnel detects the change and reloads automatically.

**Q: What happens if the SSH server goes down?**

A: SSHTunnel retries the connection every 10 seconds and re-establishes all tunnels once it reconnects.

**Q: Can I run multiple instances?**

A: Yes. Each instance needs its own config file (or a unique `--config` path) to avoid log file contention.

**Q: Does SSHTunnel support Windows as a service?**

A: There is no built-in service wrapper, but the binary runs fine under `nohup` (Linux), `start /B` (Windows), or any process supervisor.

**Q: What Go version is required to build from source?**

A: Go 1.23+ (matching `go.mod`).

**Q: Where can I report a security vulnerability?**

A: See [`SECURITY.md`](SECURITY.md) for the disclosure process.

**Q: How do I contribute?**

A: See [`CONTRIBUTING.md`](CONTRIBUTING.md).

**Q: What license is this under?**

A: MIT — see [`LICENSE`](LICENSE).

**Q: Why does `remote_ip` exist if it's always `0.0.0.0`?**

A: The field is exposed for flexibility. On servers with multiple IPs, you can pin the listener to a specific address. Most users should leave it as `"0.0.0.0"`.

**Q: Can I use `~` in the `private_key` path?**

A: No — the shell does not expand `~` inside YAML. Use an absolute path or a path relative to the working directory.

**Q: How do I stop SSHTunnel?**

A: Press `Ctrl+C`. The binary exits cleanly.

**Q: Is there a Docker image?**

A: Not yet. You can run the binary inside a container by copying the static binary into a scratch image.
