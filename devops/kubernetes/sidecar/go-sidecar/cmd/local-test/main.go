package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/example/go-sidecar/internal/k8s"
	"github.com/example/go-sidecar/internal/properties"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	propertiesFile := getEnvOrDefault("PROPERTIES_FILE", "/tmp/application.properties")
	namespace := getEnvOrDefault("POD_NAMESPACE", "default")
	configMapName := getEnvOrDefault("CONFIGMAP_NAME", "properties")

	testDir, err := os.MkdirTemp("", "sidecar-test-*")
	if err != nil {
		logger.Fatal("failed to create temp dir", zap.Error(err))
	}
	defer os.RemoveAll(testDir)

	if _, err := os.Stat(propertiesFile); os.IsNotExist(err) {
		sampleContent := "app.name=MyApp\napp.version=1.0.0\ndb.host=localhost\ndb.port=5432\ndb.name=testdb\ncache.ttl=3600\ncache.maxsize=1000\n"
		if err := os.WriteFile(propertiesFile, []byte(sampleContent), 0644); err != nil {
			logger.Fatal("failed to create sample properties file", zap.Error(err))
		}
	}

	file, err := os.Open(propertiesFile)
	if err != nil {
	}
	defer file.Close()

	allProps, err := properties.Parse(file)
	if err != nil {
		logger.Fatal("failed to parse properties file", zap.Error(err))
	}

	fmt.Println("\n✓ Sidecar successfully read properties:")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for key, value := range allProps {
		fmt.Printf("  %s = %s\n", key, value)
	}
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Printf("Total properties: %d\n", len(allProps))
	fmt.Printf("File: %s\n\n", propertiesFile)

	fmt.Println("→ Updating Kubernetes ConfigMap...")
	fmt.Printf("  Namespace: %s\n", namespace)
	fmt.Printf("  Entries: %d\n\n", len(allProps))

	k8sClient, err := k8s.NewClient(namespace, configMapName, logger)
	if err != nil {
		logger.Fatal("failed to create kubernetes client", zap.Error(err))
	}

	ctx := context.Background()
	if err := k8sClient.UpdateConfigMap(ctx, allProps); err != nil {
		logger.Fatal("failed to update ConfigMap", zap.Error(err))
	}

	fmt.Println("✓ ConfigMap updated successfully!")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
