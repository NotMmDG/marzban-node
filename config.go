package main

import (
	"os"
	"strconv"
)

// getEnv retrieves the value of the environment variable named by the key.
// It returns the value, which will be defaultValue if the variable is not present.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves the value of the environment variable named by the key
// and converts it to an integer. If conversion fails, it returns the default value.
func getEnvAsInt(name string, defaultValue int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool retrieves the value of the environment variable named by the key
// and converts it to a boolean. If conversion fails, it returns the default value.
func getEnvAsBool(name string, defaultValue bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func main() {
	servicePort := getEnvAsInt("SERVICE_PORT", 62050)
	xrayApiPort := getEnvAsInt("XRAY_API_PORT", 62051)
	xrayExecutablePath := getEnv("XRAY_EXECUTABLE_PATH", "/usr/local/bin/xray")
	xrayAssetsPath := getEnv("XRAY_ASSETS_PATH", "/usr/local/share/xray")
	sslCertFile := getEnv("SSL_CERT_FILE", "/var/lib/marzban-node/ssl_cert.pem")
	sslKeyFile := getEnv("SSL_KEY_FILE", "/var/lib/marzban-node/ssl_key.pem")
	sslClientCertFile := getEnv("SSL_CLIENT_CERT_FILE", "")
	debug := getEnvAsBool("DEBUG", false)
	serviceProtocol := getEnv("SERVICE_PROTOCOL", "rpyc")

	// Use the variables as needed in your application
}
