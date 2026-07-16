package sshlib

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

// --- helpers ---

func generateTestKey(t *testing.T) (ssh.Signer, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	signer, err := ssh.ParsePrivateKey(privPEM)
	if err != nil {
		t.Fatalf("failed to parse private key: %v", err)
	}
	return signer, string(privPEM)
}

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
	return path
}

func startTestSSHServer(t *testing.T, signer ssh.Signer) (string, func()) {
	t.Helper()
	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}
	config.AddHostKey(signer)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			return
		}
		go ssh.DiscardRequests(reqs)
		go func() {
			for newCh := range chans {
				ch, _, err := newCh.Accept()
				if err != nil {
					continue
				}
				go func(c ssh.Channel) {
					defer c.Close()
					io.Copy(io.Discard, c)
				}(ch)
			}
		}()
		sshConn.Wait()
	}()

	addr := listener.Addr().String()
	cleanup := func() {
		listener.Close()
		wg.Wait()
	}
	return addr, cleanup
}

func newTestConfig(host, keyPath string, tunnels []Tunnel) *Config {
	if tunnels == nil {
		tunnels = []Tunnel{
			{LocalIP: "127.0.0.1", LocalPort: 9999, RemoteIP: "0.0.0.0", RemotePort: 19999},
		}
	}
	return &Config{
		SSHConfig: SSHConfig{
			Host:       host,
			Port:       22,
			User:       "test",
			PrivateKey: keyPath,
			Tunnels:    tunnels,
		},
	}
}

func captureBanner() func() string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	return func() string {
		w.Close()
		os.Stdout = old
		var buf bytes.Buffer
		io.Copy(&buf, r)
		return buf.String()
	}
}

// --- tests ---

func TestPrintBanner(t *testing.T) {
	done := captureBanner()
	PrintBanner()
	out := done()
	if !strings.Contains(out, "SSH Tunneling Service") {
		t.Errorf("banner missing title: %s", out)
	}
}

func TestLoadConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	yamlContent := `
ssh_config:
  host: "example.com"
  port: 2222
  user: "alice"
  private_key: "/home/alice/.ssh/id_rsa"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 8080
      remote_ip: "0.0.0.0"
      remote_port: 9090
`
	path := writeTempFile(t, dir, "sshtunnel.yml", yamlContent)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.SSHConfig.Host != "example.com" {
		t.Errorf("host = %q, want %q", cfg.SSHConfig.Host, "example.com")
	}
	if cfg.SSHConfig.Port != 2222 {
		t.Errorf("port = %d, want %d", cfg.SSHConfig.Port, 2222)
	}
	if cfg.SSHConfig.User != "alice" {
		t.Errorf("user = %q, want %q", cfg.SSHConfig.User, "alice")
	}
	if len(cfg.SSHConfig.Tunnels) != 1 {
		t.Fatalf("len(tunnels) = %d, want 1", len(cfg.SSHConfig.Tunnels))
	}
	if cfg.SSHConfig.Tunnels[0].LocalPort != 8080 {
		t.Errorf("tunnel local_port = %d, want 8080", cfg.SSHConfig.Tunnels[0].LocalPort)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "does_not_exist.yml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.yml", `ssh_config: {invalid`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestLoadConfig_EmptyTunnels(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "empty.yml", `
ssh_config:
  host: "example.com"
  port: 22
  user: "test"
  private_key: "/tmp/key"
  tunnels: []
`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if len(cfg.SSHConfig.Tunnels) != 0 {
		t.Errorf("expected 0 tunnels, got %d", len(cfg.SSHConfig.Tunnels))
	}
}

func TestLoadConfig_MultipleTunnels(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "multi.yml", `
ssh_config:
  host: "example.com"
  port: 22
  user: "test"
  private_key: "/tmp/key"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 80
      remote_port: 8080
    - local_ip: "127.0.0.1"
      local_port: 5432
      remote_port: 15432
    - local_ip: "0.0.0.0"
      local_port: 3000
      remote_port: 3000
`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if len(cfg.SSHConfig.Tunnels) != 3 {
		t.Fatalf("expected 3 tunnels, got %d", len(cfg.SSHConfig.Tunnels))
	}
	ports := []int{80, 5432, 3000}
	for i, tr := range cfg.SSHConfig.Tunnels {
		if tr.LocalPort != ports[i] {
			t.Errorf("tunnel[%d] local_port = %d, want %d", i, tr.LocalPort, ports[i])
		}
	}
}

func TestLoadConfig_DefaultPort(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "noport.yml", `
ssh_config:
  host: "example.com"
  user: "test"
  private_key: "/tmp/key"
  tunnels:
    - local_port: 80
      remote_port: 8080
`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.SSHConfig.Port != 0 {
		t.Errorf("port = %d, want 0 (unset)", cfg.SSHConfig.Port)
	}
}

func TestConfigYAMLRoundTrip(t *testing.T) {
	original := Config{
		SSHConfig: SSHConfig{
			Host:       "vps.example.com",
			Port:       2222,
			User:       "deploy",
			PrivateKey: "/home/deploy/.ssh/id_ed25519",
			Tunnels: []Tunnel{
				{LocalIP: "0.0.0.0", LocalPort: 8080, RemoteIP: "0.0.0.0", RemotePort: 9090},
				{LocalIP: "127.0.0.1", LocalPort: 3000, RemoteIP: "0.0.0.0", RemotePort: 3000},
			},
		},
	}
	data, err := yaml.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var decoded Config
	if err := yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if decoded.SSHConfig.Host != original.SSHConfig.Host {
		t.Errorf("host mismatch: %q vs %q", decoded.SSHConfig.Host, original.SSHConfig.Host)
	}
	if len(decoded.SSHConfig.Tunnels) != len(original.SSHConfig.Tunnels) {
		t.Errorf("tunnel count mismatch: %d vs %d", len(decoded.SSHConfig.Tunnels), len(original.SSHConfig.Tunnels))
	}
}

func TestSSHDial_InvalidKeyPath(t *testing.T) {
	cfg := &SSHConfig{
		Host:       "127.0.0.1",
		Port:       22,
		User:       "test",
		PrivateKey: "/nonexistent/key",
	}
	_, err := sshDial(cfg)
	if err == nil {
		t.Fatal("expected error for invalid key path, got nil")
	}
}

func TestSSHDial_InvalidKeyContent(t *testing.T) {
	dir := t.TempDir()
	badKey := writeTempFile(t, dir, "bad_key", "not a valid private key")
	cfg := &SSHConfig{
		Host:       "127.0.0.1",
		Port:       22,
		User:       "test",
		PrivateKey: badKey,
	}
	_, err := sshDial(cfg)
	if err == nil {
		t.Fatal("expected error for invalid key content, got nil")
	}
}

func TestSSHDial_ConnectionRefused(t *testing.T) {
	_, keyPEM := generateTestKey(t)
	dir := t.TempDir()
	keyPath := writeTempFile(t, dir, "valid_key", keyPEM)
	cfg := &SSHConfig{
		Host:       "127.0.0.1",
		Port:       1,
		User:       "test",
		PrivateKey: keyPath,
	}
	_, err := sshDial(cfg)
	if err == nil {
		t.Fatal("expected connection error, got nil")
	}
}

func TestSSHDial_ValidKey(t *testing.T) {
	signer, keyPEM := generateTestKey(t)
	addr, cleanup := startTestSSHServer(t, signer)
	defer cleanup()

	dir := t.TempDir()
	keyPath := writeTempFile(t, dir, "id_rsa", keyPEM)

	cfg := &SSHConfig{
		Host:       strings.Split(addr, ":")[0],
		Port:       mustPort(t, addr),
		User:       "test",
		PrivateKey: keyPath,
	}
	client, err := sshDial(cfg)
	if err != nil {
		t.Fatalf("sshDial failed: %v", err)
	}
	client.Close()
}

func TestSetupLogging(t *testing.T) {
	dir, err := os.MkdirTemp("", "sshtunneltest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	prevOutput := log.Writer()
	SetupLogging(dir)
	log.Println("test log entry")
	log.SetOutput(prevOutput)

	logPath := filepath.Join(dir, "ssh_tunneling.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(data), "test log entry") {
		t.Errorf("log file missing entry: %s", data)
	}
}

func TestSetupLogging_OutputToFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "sshtunneltest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	prevOutput := log.Writer()
	SetupLogging(dir)
	log.Println("file output test")
	log.SetOutput(prevOutput)

	logPath := filepath.Join(dir, "ssh_tunneling.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	if !strings.Contains(string(data), "file output test") {
		t.Errorf("log file missing output: %s", data)
	}
}

func TestHandleTunnel_LocalRefused(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	done := make(chan struct{})
	go func() {
		handleTunnel(client, "127.0.0.1", 1)
		close(done)
	}()

	server.Write([]byte("hello"))
	server.Close()
	<-done

	if !strings.Contains(buf.String(), "Failed to connect to local service") {
		t.Errorf("expected connection refused log, got: %s", buf.String())
	}
}

type pipeConn struct {
	io.ReadCloser
	io.WriteCloser
}

func (p *pipeConn) Close() error {
	e1 := p.ReadCloser.Close()
	e2 := p.WriteCloser.Close()
	if e1 != nil {
		return e1
	}
	return e2
}

func (p *pipeConn) LocalAddr() net.Addr                { return fakeAddr{} }
func (p *pipeConn) RemoteAddr() net.Addr               { return fakeAddr{} }
func (p *pipeConn) SetDeadline(t time.Time) error      { return nil }
func (p *pipeConn) SetReadDeadline(t time.Time) error  { return nil }
func (p *pipeConn) SetWriteDeadline(t time.Time) error { return nil }

type fakeAddr struct{}

func (fakeAddr) Network() string { return "pipe" }
func (fakeAddr) String() string  { return "pipe" }

func TestHandleTunnel_Success(t *testing.T) {
	serverReady := make(chan struct{})
	received := make(chan []byte, 1)

	tcpListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	localPort := tcpListener.Addr().(*net.TCPAddr).Port

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conn, err := tcpListener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		close(serverReady)
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if n > 0 {
			received <- buf[:n]
		}
		// Wait so handleTunnel can attempt io.Copy(conn, localConn)
		// without the TCP connection closing immediately.
		time.Sleep(100 * time.Millisecond)
	}()
	defer func() {
		tcpListener.Close()
		wg.Wait()
	}()

	tunR, tunW := io.Pipe()
	appR, appW := io.Pipe()

	client := &pipeConn{ReadCloser: appR, WriteCloser: tunW}
	server := &pipeConn{ReadCloser: tunR, WriteCloser: appW}

	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)

	done := make(chan struct{})
	go func() {
		handleTunnel(client, "127.0.0.1", localPort)
		close(done)
	}()

	<-serverReady
	server.Write([]byte("hello"))
	select {
	case data := <-received:
		if string(data) != "hello" {
			t.Errorf("got %q, want %q", string(data), "hello")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for data on local service")
	}

	server.Close()
	<-done

	if !strings.Contains(logBuf.String(), "Connection established to local service") {
		t.Errorf("expected connection log, got: %s", logBuf.String())
	}
}

func TestStartTunnel_InvalidRemote(t *testing.T) {
	signer, keyPEM := generateTestKey(t)
	addr, cleanup := startTestSSHServer(t, signer)
	defer cleanup()

	dir := t.TempDir()
	keyPath := writeTempFile(t, dir, "id_rsa", keyPEM)

	sshCfg := &SSHConfig{
		Host:       strings.Split(addr, ":")[0],
		Port:       mustPort(t, addr),
		User:       "test",
		PrivateKey: keyPath,
	}
	client, err := sshDial(sshCfg)
	if err != nil {
		t.Fatalf("sshDial failed: %v", err)
	}
	defer client.Close()

	tun := Tunnel{LocalIP: "127.0.0.1", LocalPort: 9999, RemoteIP: "0.0.0.0", RemotePort: 0}
	err = startTunnel(client, sshCfg, tun)
	if err == nil {
		t.Fatal("expected error for invalid remote port, got nil")
	}
}

func TestMaintainSSHConnection_BadConfig(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	cfg := &Config{
		SSHConfig: SSHConfig{
			Host:       "127.0.0.1",
			Port:       1,
			User:       "test",
			PrivateKey: "/nonexistent",
			Tunnels:    []Tunnel{{LocalIP: "127.0.0.1", LocalPort: 9999, RemoteIP: "0.0.0.0", RemotePort: 19999}},
		},
	}

	done := make(chan struct{})
	go func() {
		MaintainSSHConnection(cfg)
		close(done)
	}()

	time.Sleep(2 * time.Second)
	out := buf.String()
	if !strings.Contains(out, "Failed to establish SSH connection") {
		t.Errorf("expected retry log, got: %s", out)
	}
}

// --- benchmarks ---

func BenchmarkLoadConfig(b *testing.B) {
	dir := b.TempDir()
	content := `
ssh_config:
  host: "bench.example.com"
  port: 22
  user: "bench"
  private_key: "/bench/key"
  tunnels:
    - local_ip: "0.0.0.0"
      local_port: 80
      remote_port: 8080
`
	path := filepath.Join(dir, "bench.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		b.Fatalf("failed to write bench.yml: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadConfig(path)
	}
}

// --- helpers for tests ---

func mustPort(t *testing.T, addr string) int {
	t.Helper()
	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("SplitHostPort(%q) failed: %v", addr, err)
	}
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		t.Fatalf("parse port %q: %v", portStr, err)
	}
	return port
}
