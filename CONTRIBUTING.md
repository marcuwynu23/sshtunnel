# Contributing to SSHTunnel

Thank you for your interest in contributing! Whether it's a bug report, feature request, or pull request, we appreciate your help.

- [Code of Conduct](#code-of-conduct)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Makefile Reference](#makefile-reference)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Conventions](#commit-conventions)
- [PR Process](#pr-process)
- [Release Process](#release-process)
- [Questions](#questions)

---

## Code of Conduct

This project is governed by the [Contributor Covenant](CODE_OF_CONDUCT.md). By participating, you agree to uphold its standards.

---

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.23+ | Compiler and standard toolchain |
| make | any | Build automation (optional — `go build` works directly) |
| git | any | Version control |

---

## Project Structure

```
sshtunnel/
├── main.go                 # Entry point: CLI parsing, config loading, file watcher
├── sshlib/
│   └── tunnel.go           # SSH dial, reverse tunnel loop, logging, banner
├── sshtunnel.yml           # Default configuration file (packaged with releases)
├── makefile                # Cross-compilation build targets
├── resources/
│   └── logo.svg            # Project logo for README
├── .github/
│   ├── workflows/
│   │   ├── test.yml        # CI: lint, vet, test, build on push/PR to main
│   │   └── release.yml     # CD: cross-compile + publish on tag push
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.yml
│   │   ├── feature_request.yml
│   │   └── config.yml
│   └── PULL_REQUEST_TEMPLATE.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── USER-GUIDE.md
├── LICENSE
└── README.md
```

---

## Makefile Reference

| Target | Description |
|---|---|
| `dev` | Build for the local platform and copy to `D:\Executables\sshtunnel` (Windows) |
| `ci` | Run `go fmt`, `go vet`, `go test`, `go build` (CI entry point) |
| `all` | Build and package all supported OS/arch combinations |
| `clean` | Remove the `build/` directory |
| `build` | Run `clean` then build all targets sequentially |

Platform-specific targets (each produces a `.tar.gz` in `build/`):

| Target | OS | Arch |
|---|---|---|
| `linux_amd64` | Linux | amd64 |
| `linux_386` | Linux | 386 |
| `linux_arm64` | Linux | arm64 |
| `linux_arm` | Linux | arm |
| `windows_amd64` | Windows | amd64 |
| `windows_386` | Windows | 386 |
| `macos_amd64` | macOS | amd64 |

Example:

```bash
make linux_amd64 VERSION=v1.2.3
```

---

## Development Workflow

1. **Fork** the repository and clone it locally.
2. **Create a branch** with a descriptive name:

   ```bash
   git checkout -b fix/reconnect-backoff
   ```

3. **Make your changes.** Follow the [coding standards](#coding-stanards) below.
4. **Run the CI checks locally:**

   ```bash
   make ci
   ```

5. **Commit** using [conventional commits](#commit-conventions).
6. **Push** and open a pull request against `main`.

---

## Coding Standards

- **Formatting:** Run `go fmt ./...` before committing.
- **Imports:** Group standard library, third-party, and local imports with a blank line between each group.
- **Naming:** Follow Go conventions — `camelCase` for unexported, `PascalCase` for exported. Acronyms are all-uppercase (`SSHConfig`, not `SshConfig`).
- **Errors:** Always check and handle errors. Prefer `fmt.Errorf("context: %w", err)` for wrapping.
- **Logging:** Use the standard `log` package (not `fmt.Print` for diagnostics).
- **No unused code:** Run `go vet ./...` to detect dead code.

---

## Testing

- Run tests with:

  ```bash
  go test ./...
  ```

- Tests should be **table-driven** where possible.
- New functionality must include tests that cover the happy path and at least one error case.
- The CI workflow (`make ci`) runs `go test` on every push and pull request.

---

## Commit Conventions

This project uses **Conventional Commits**. The tool itself generates release notes from commit history, so good messages matter.

```
<type>(<scope>): <description>
```

| Type | Usage |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation only changes |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `test` | Adding or correcting tests |
| `ci` | CI/CD configuration or script changes |
| `chore` | Build process, tooling, or dependency changes |

### Scope

The scope should be the package or area affected. Example scopes: `ssh`, `config`, `watcher`, `build`, `docs`.

### Examples

```
feat(ssh): add exponential backoff to reconnect loop
fix(config): handle YAML with empty tunnels array
docs(readme): add CI/CD integration section
ci(release): fix artifact naming for Windows builds
```

Breaking changes should include `BREAKING CHANGE:` in the footer:

```
feat(config): switch from YAML v1 to v2 schema

BREAKING CHANGE: the `tunnels` field now requires `remote_ip`.
```

---

## PR Process

1. Base your PR on the latest `main` branch.
2. Ensure the title follows conventional commits (e.g., `feat(ssh): add keep-alive`).
3. Fill out the [pull request template](.github/PULL_REQUEST_TEMPLATE.md).
4. A maintainer will review your PR. They may request changes before merging.
5. CI must pass before merge.

### Checklist before submitting

- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes
- [ ] `go vet ./...` is clean
- [ ] I have read `CONTRIBUTING.md`
- [ ] Docs are updated if the change affects the public API, config schema, or CLI flags

---

## Release Process

1. A maintainer pushes a tag matching `v*` (e.g., `v1.2.0`).
2. The [Release workflow](.github/workflows/release.yml) runs automatically:
   - Cross-compiles all targets via `make build VERSION=<tag>`
   - Packages each binary with `sshtunnel.yml` into `.tar.gz`
   - Creates a GitHub Release with the archives
3. The changelog is maintained externally via `auto-changelog` (see `CHANGELOG.md`).

---

## Questions

If you have questions or want to discuss an idea before writing code:

- Open a [Discussion](https://github.com/marcuwynu23/sshtunnel/discussions)
- Open an [Issue](https://github.com/marcuwynu23/sshtunnel/issues)

For security vulnerabilities, see [`SECURITY.md`](SECURITY.md).
