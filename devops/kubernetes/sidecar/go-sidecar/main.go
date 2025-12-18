package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/example/go-sidecar/internal/config"
	"github.com/example/go-sidecar/internal/k8s"
	"github.com/example/go-sidecar/internal/properties"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	logger.Info("starting sidecar",
		zap.String("properties_file", cfg.PropertiesFile),
		zap.String("namespace", cfg.Namespace),
		zap.String("configmap_name", cfg.ConfigMapName),
	)

	k8sClient, err := k8s.NewClient(cfg.Namespace, cfg.ConfigMapName, logger)
	if err != nil {
		logger.Fatal("failed to create kubernetes client", zap.Error(err))
	}

	props, err := readPropertiesFile(cfg.PropertiesFile, logger)
	if err != nil {
		logger.Fatal("failed to read properties file", zap.Error(err))
	}

	ctx := context.Background()
	if err := k8sClient.UpdateConfigMap(ctx, props); err != nil {
		logger.Fatal("failed to update ConfigMap", zap.Error(err))
	}

	logger.Info("sidecar completed successfully, exiting")
}

func readPropertiesFile(filePath string, logger *zap.Logger) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open properties file: %w", err)
	}
	defer file.Close()

	props, err := properties.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse properties file: %w", err)
	}

	logger.Info("loaded properties", zap.String("file", filePath), zap.Int("total", len(props)))
	return props, nil
}
