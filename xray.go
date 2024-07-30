package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"
)

type XRayConfig struct {
	Config map[string]interface{}
	ApiPort    int
	SSLCert    string
	SSLKey     string
	PeerIP     string
}

type XRayCore struct {
	ExecutablePath string
	AssetsPath     string
	Version        string
	Process        *exec.Cmd
	Restarting     bool
	LogsBuffer     *deque.Deque
	OnStartFuncs   []func()
	OnStopFuncs    []func()
	Env            map[string]string
	mu             sync.Mutex
}

func NewXRayConfig(configStr, peerIP string) (*XRayConfig, error) {
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return nil, err
	}

	xrayConfig := &XRayConfig{
		Config: config,
		ApiPort: XRAY_API_PORT, // Define XRAY_API_PORT appropriately
		SSLCert: SSL_CERT_FILE, // Define SSL_CERT_FILE appropriately
		SSLKey:  SSL_KEY_FILE,  // Define SSL_KEY_FILE appropriately
		PeerIP:  peerIP,
	}
	xrayConfig.applyAPI()
	return xrayConfig, nil
}

func (x *XRayConfig) applyAPI() {
	// Remove existing 'dokodemo-door' inbounds and API rules
	// Apply API configuration as done in the original script
	// This requires converting the relevant Python dictionary operations to Go map operations
}

func (x *XRayConfig) ToJSON() (string, error) {
	bytes, err := json.Marshal(x.Config)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewXRayCore(executablePath, assetsPath string) *XRayCore {
	return &XRayCore{
		ExecutablePath: executablePath,
		AssetsPath:     assetsPath,
		LogsBuffer:     deque.New(),
		Env: map[string]string{
			"XRAY_LOCATION_ASSET": assetsPath,
		},
	}
}

func (x *XRayCore) GetVersion() (string, error) {
	cmd := exec.Command(x.ExecutablePath, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`^Xray (\d+\.\d+\.\d+)`)
	match := re.FindStringSubmatch(string(output))
	if len(match) > 1 {
		return match[1], nil
	}
	return "", nil
}

func (x *XRayCore) captureProcessLogs(debug bool) {
	// Function to capture logs, similar to the Python version
	// Start capturing stdout and stderr from the process
}

func (x *XRayCore) Start(config *XRayConfig) error {
	if x.Process != nil {
		return fmt.Errorf("Xray is already started")
	}

	configJSON, err := config.ToJSON()
	if err != nil {
		return err
	}

	cmd := exec.Command(x.ExecutablePath, "run", "-config", "stdin:")
	cmd.Env = os.Environ()
	for k, v := range x.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stdin = strings.NewReader(configJSON)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	x.Process = cmd

	go func() {
		if err := cmd.Run(); err != nil {
			log.Printf("Process finished with error: %v", err)
		}
	}()

	x.captureProcessLogs(DEBUG) // Define DEBUG appropriately

	for _, fn := range x.OnStartFuncs {
		go fn()
	}

	return nil
}

func (x *XRayCore) Stop() {
	if x.Process != nil {
		x.Process.Process.Signal(syscall.SIGTERM)
		x.Process.Wait()
		x.Process = nil
	}

	for _, fn := range x.OnStopFuncs {
		go fn()
	}
}

func (x *XRayCore) Restart(config *XRayConfig) {
	x.mu.Lock()
	defer x.mu.Unlock()

	if x.Restarting {
		return
	}

	x.Restarting = true
	log.Println("Restarting Xray core...")

	x.Stop()
	if err := x.Start(config); err != nil {
		log.Printf("Error restarting Xray: %v", err)
	}

	x.Restarting = false
}

func (x *XRayCore) OnStart(fn func()) {
	x.OnStartFuncs = append(x.OnStartFuncs, fn)
}

func (x *XRayCore) OnStop(fn func()) {
	x.OnStopFuncs = append(x.OnStopFuncs, fn)
}

func main() {
	// Define your config, setup signals, and run the core.
	// For example:
	/*
		configStr := `{"inbounds": [...], "outbounds": [...]}`
		peerIP := "127.0.0.1"
		config, err := NewXRayConfig(configStr, peerIP)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		xrayCore := NewXRayCore("/usr/bin/xray", "/usr/share/xray")
		if err := xrayCore.Start(config); err != nil {
			log.Fatalf("Failed to start Xray: %v", err)
		}
	*/
	// Listen for system signals and handle graceful shutdown, etc.
}
