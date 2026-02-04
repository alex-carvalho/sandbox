package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	PropertiesFile string
	WatchDir       string
	Namespace      string
	ConfigMapName  string
}

func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		PropertiesFile: getEnvOrDefault("PROPERTIES_FILE", "/etc/config/application.properties"),
		WatchDir:       getEnvOrDefault("WATCH_DIR", "/etc/config"),
		Namespace:      getNamespace(),
		ConfigMapName:  getEnvOrDefault("CONFIGMAP_NAME", "properties"),
	}

	if cfg.PropertiesFile == "" {
		return nil, fmt.Errorf("PROPERTIES_FILE is required")
	}

	if cfg.Namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	if cfg.ConfigMapName == "" {
		return nil, fmt.Errorf("configmap name is required")
	}

	return cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// getNamespace returns the namespace from environment variable or from the service account token
func getNamespace() string {
	if ns, ok := os.LookupEnv("POD_NAMESPACE"); ok && ns != "" {
		return ns
	}

	nsBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		ns := strings.TrimSpace(string(nsBytes))
		if ns != "" {
			return ns
		}
	}

	// Fallback to default namespace
	return "default"
}
