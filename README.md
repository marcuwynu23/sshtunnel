<div align="center">

<img src="./resources/logo.svg" width="420" alt="SSHTunnel logo">

<a href="https://github.com/marcuwynu23/sshtunnel/releases"><img src="https://img.shields.io/github/v/release/marcuwynu23/sshtunnel" alt="Release"></a>
<a href="https://github.com/marcuwynu23/sshtunnel/blob/main/LICENSE"><img src="https://img.shields.io/github/license/marcuwynu23/sshtunnel?logo=github" alt="License"></a>
<a href="https://github.com/marcuwynu23/sshtunnel"><img src="https://img.shields.io/github/stars/marcuwynu23/sshtunnel" alt="Stars"></a>
<img src="https://img.shields.io/github/go-mod/go-version/marcuwynu23/sshtunnel" alt="Go version">

<strong>Reverse SSH tunnels, configured in YAML.</strong> A cross-platform CLI tool that reads a simple config file and keeps your tunnels alive — automatically reconnecting on failure and hot-reloading on config changes.

➡️ **[Read the full user guide →](USER-GUIDE.md)**

</div>

---

## Table of Contents

- [What Is SSHTunnel?](#what-is-sshtunnel)
- [Use Cases](#use-cases)
- [Benefits](#benefits-for-developers)
- [Comparison](#advantages-over-other-tools)
- [User Guide](USER-GUIDE.md)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Reference](#cli-commands)
- [Configuration](#configuration)
- [Development](#development)
- [Architecture](#architecture)

---

## What Is SSHTunnel?

**SSHTunnel** is a Go-based CLI tool that establishes and maintains **reverse SSH tunnels** from a single YAML configuration file. It handles reconnection, hot-reloads when the config changes, and logs to both console and file.

### What It Does

- **Configure in YAML** — Define SSH host, port, user, private key, and one or more tunnels in a single file
- **Reverse-tunnel automatically** — Each tunnel forwards a remote port back to a local service
- **Reconnect on failure** — Drops and retries every 10 seconds if the SSH connection dies
- **Hot-reload config** — Edit `sshtunnel.yml` while it runs; tunnels restart automatically
- **Log to file + console** — Writes structured logs to `ssh_tunneling.log` alongside the config file
- **Cross-compile to 8 targets** — Linux (amd64/386/arm64/arm), Windows (amd64/386), macOS (amd64)
- **Zero runtime dependencies** — Single static binary, no SSH client required on the target

### Why Use It?

| Problem                                         | How SSHTunnel Solves It                                                                                                      |
| ----------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Services behind NAT or firewall are unreachable | **Reverse tunnel** — the machine behind NAT initiates the outbound SSH connection and exposes local ports on a public server |
| SSH commands are long and easy to mistype       | **YAML config** — declare tunnels once in a structured file; running is just `sshtunnel`                                     |
| SSH connections drop unpredictably              | **Auto-reconnect** — retries the connection every 10 seconds and re-establishes all tunnels                                  |
| Config changes require restarting tunnels       | **Hot-reload** — edit `sshtunnel.yml` and the tool picks it up via `fsnotify`                                                |
| Hard to tell what happened while you were away  | **Dual logging** — every event goes to both the terminal and `ssh_tunneling.log`                                             |

### The Philosophy

1. **Configuration over invocation.** You shouldn't have to type 200 characters of SSH flags every time. Define everything in YAML once and run a single command.
2. **Resilient by default.** Connections fail. SSHTunnel assumes they will and bounces back automatically — no `systemd` unit or supervisor process required.
3. **Your data stays yours.** No cloud, no telemetry, no accounts. The binary and a config file are all you need.

---

## Use Cases

| Scenario                                       | How SSHTunnel Helps                                                                       |
| ---------------------------------------------- | ----------------------------------------------------------------------------------------- |
| **Expose a local dev server to the internet**  | Forward `localhost:8080` → `public-server:9090` for testing webhooks                      |
| **Remote Raspberry Pi management**             | Expose SSH (port 22) behind a home NAT to a VPS on port 2222                              |
| **Database access from a locked-down network** | Tunnel a local `postgres:5432` to a remote port so your office can reach it               |
| **CI/CD runner behind a firewall**             | Self-hosted runners can expose build metrics or debug ports through a single outbound SSH |

---

## Benefits for Developers

- **~10-second setup** — Write 5 lines of YAML and run the binary
- **Single binary** — No runtime, no dependencies, no package manager required
- **Multi-arch builds** — `make all` produces tarballs for Linux, Windows, and macOS on x86 and ARM
- **Config hot-reload** — No downtime when adding or removing tunnels
- **Automatic reconnection** — Survives network blips and server restarts
- **Simple CLI** — Two flags (`--config` and `--help`), no subcommands
- **MIT licensed** — Free for personal and commercial use
- **Private key authentication** — Uses your existing SSH keys; no passwords stored in config

---

## Advantages Over Other Tools

| Aspect                  | SSHTunnel             | `autossh`              | `ssh -R` manual        | `ngrok` / `bore`         |
| ----------------------- | --------------------- | ---------------------- | ---------------------- | ------------------------ |
| **Setup time**          | ~10 seconds           | ~2 minutes             | Repeating effort       | ~30 seconds              |
| **Config file**         | YAML (single file)    | Environment / args     | None (CLI flags)       | CLI flags / TOML         |
| **Hot-reload**          | Yes (fsnotify)        | No                     | No                     | Partial (CLI restart)    |
| **Auto-reconnect**      | Built-in (10 s retry) | Built-in (monitors)    | No                     | Built-in                 |
| **Multi-tunnel**        | Array in YAML         | Multiple processes     | Multiple commands      | Per-tunnel process       |
| **Auth method**         | Private key           | Key / password / agent | Key / password / agent | Service account          |
| **Cross-platform**      | Linux, Windows, macOS | Linux, macOS           | Everywhere with SSH    | Linux, Windows, macOS    |
| **Third-party service** | No (your own server)  | No                     | No                     | Yes (ngrok/bore servers) |
| **License**             | MIT                   | Custom (GPL-like)      | OpenSSH license        | Proprietary / Apache 2.0 |
| **Binary size**         | ~5 MB (compressed)    | ~100 KB (+ SSH)        | Built-in               | ~10 MB                   |

---

## Installation

### Download a Prebuilt Binary

Grab the latest tarball for your platform from the [releases page](https://github.com/marcuwynu23/sshtunnel/releases):

| Platform      | File                                       |
| ------------- | ------------------------------------------ |
| Linux amd64   | `sshtunnel_linux_amd64_<version>.tar.gz`   |
| Linux 386     | `sshtunnel_linux_386_<version>.tar.gz`     |
| Linux arm64   | `sshtunnel_linux_arm64_<version>.tar.gz`   |
| Linux arm     | `sshtunnel_linux_arm_<version>.tar.gz`     |
| Windows amd64 | `sshtunnel_windows_amd64_<version>.tar.gz` |
| Windows 386   | `sshtunnel_windows_386_<version>.tar.gz`   |
| macOS amd64   | `sshtunnel_macos_amd64_<version>.tar.gz`   |

Extract and run:

```bash
tar xzf sshtunnel_linux_amd64_<version>.tar.gz
cd sshtunnel_linux_amd64_<version>
./sshtunnel --help
```

### Build from Source

```bash
git clone https://github.com/marcuwynu23/sshtunnel.git
cd sshtunnel
go build -o sshtunnel main.go
./sshtunnel --help
```

Requires **Go 1.23+**.

---

## Quick Start

Create a file called `sshtunnel.yml` in the same directory as the binary:

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

Then run:

```bash
sshtunnel
```

The tool connects to `203.0.113.1`, opens a reverse tunnel so that port 9090 on the remote server forwards traffic to port 8080 on your machine, and begins logging to `ssh_tunneling.log`.

---

## CLI Commands

### `sshtunnel`

Start the SSH tunneling service using the config file found next to the binary.

```bash
sshtunnel [--config <path>]
```

| Flag       | Default                          | Description                         |
| ---------- | -------------------------------- | ----------------------------------- |
| `--config` | `<executable_dir>/sshtunnel.yml` | Path to the YAML configuration file |
| `--help`   | —                                | Show usage text and exit            |

#### Examples

**Use the default config (next to the binary):**

```bash
sshtunnel
```

**Use a specific config file:**

```bash
sshtunnel --config ./production.yml
```

**Show help:**

```bash
sshtunnel --help
```

---

## Configuration

The entire tool is driven by a single YAML file (`sshtunnel.yml`):

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
    - local_ip: "127.0.0.1"
      local_port: 3000
      remote_ip: "0.0.0.0"
      remote_port: 3000
```

### Fields

| Key                      | Type   | Default     | Description                                                                                   |
| ------------------------ | ------ | ----------- | --------------------------------------------------------------------------------------------- |
| `ssh_config.host`        | string | —           | Remote SSH server hostname or IP                                                              |
| `ssh_config.port`        | int    | `22`        | SSH port on the remote server                                                                 |
| `ssh_config.user`        | string | —           | SSH login username                                                                            |
| `ssh_config.private_key` | string | —           | Path to the private key file for authentication                                               |
| `ssh_config.tunnels`     | array  | `[]`        | List of tunnel definitions                                                                    |
| `tunnels[].local_ip`     | string | `"0.0.0.0"` | Local IP to bind the incoming side of the tunnel                                              |
| `tunnels[].local_port`   | int    | —           | Local port to forward traffic to                                                              |
| `tunnels[].remote_ip`    | string | `"0.0.0.0"` | Remote IP to bind the listener on the SSH server (must be `0.0.0.0` or an IP the server owns) |
| `tunnels[].remote_port`  | int    | —           | Remote port to expose on the SSH server                                                       |

### Config Precedence

The `--config` CLI flag overrides the default path (`<executable_dir>/sshtunnel.yml`). There is no cascade from environment variables or system-wide paths.

---

## Development

### Prerequisites

| Tool | Version | Purpose                                                       |
| ---- | ------- | ------------------------------------------------------------- |
| Go   | 1.23+   | Compiler and toolchain                                        |
| make | any     | Build automation (optional — you can run `go build` directly) |

### Commands

```bash
git clone https://github.com/marcuwynu23/sshtunnel.git
cd sshtunnel

# Build for your current platform
make dev          # builds sshtunnel.exe and copies it to D:\Executables\sshtunnel (Windows)

# Full CI check (fmt → vet → test → build)
make ci

# Cross-compile for all supported platforms
make all

# Run the tool locally
./sshtunnel --config sshtunnel.yml
```

Every `make` target is documented in the [Makefile Reference](USER-GUIDE.md#makefile-reference) section of the user guide.

### Project Structure

```
sshtunnel/
├── main.go              # Entry point — CLI, config load, file watcher
├── sshlib/
│   └── tunnel.go        # SSH dial, tunnel loop, log setup, banner
├── sshtunnel.yml        # Example / default config file
├── makefile             # Cross-compilation build system
├── resources/
│   └── logo.svg         # Project logo
├── .github/
│   ├── workflows/
│   │   ├── test.yml     # CI: run on push/PR to main
│   │   └── release.yml  # CD: build & publish on tag push
│   ├── ISSUE_TEMPLATE/  # Bug report and feature request templates
│   └── PULL_REQUEST_TEMPLATE.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
└── LICENSE
```

---

## Architecture

SSHTunnel is a single-process Go application with three concurrent concerns:

1. **CLI / Config** — `main.go` parses flags, resolves the config path, loads the YAML, and passes it to the tunnel manager.
2. **SSH Tunnel Manager** — `sshlib/tunnel.go` dials the SSH server, opens a reverse listener for each tunnel definition, and copies traffic bidirectionally. If the connection drops, it retries in a 10-second loop.
3. **File Watcher** — `main.go` uses `fsnotify` to monitor the config file for `Write` events. On change, it re-reads the YAML and re-initializes the tunnel manager — no restart required.

All components share the standard `log` package, which is wired to a `io.MultiWriter` that writes to both `os.Stdout` and `ssh_tunneling.log`.

---

## Contributing

We welcome contributions! Read [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full guide on setting up a development environment, coding standards, commit conventions, and the PR process.

## License

SSHTunnel is licensed under the MIT License. See [`LICENSE`](LICENSE) for details.
