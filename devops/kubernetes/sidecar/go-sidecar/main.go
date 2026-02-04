package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/example/go-sidecar/internal/config"
	"github.com/example/go-sidecar/internal/k8s"
	"github.com/example/go-sidecar/internal/watcher"
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
		zap.String("watch_dir", cfg.WatchDir),
		zap.String("namespace", cfg.Namespace),
		zap.String("configmap_name", cfg.ConfigMapName),
	)

	k8sClient, err := k8s.NewClient(cfg.Namespace, cfg.ConfigMapName, logger)
	if err != nil {
		logger.Fatal("failed to create kubernetes client", zap.Error(err))
	}

	w, err := watcher.NewWatcher(cfg.WatchDir, k8sClient, logger)
	if err != nil {
		logger.Fatal("failed to create watcher", zap.Error(err))
	}
	defer w.Close()

	logger.Info("sidecar watcher started, watching for changes")
	ctx := context.Background()
	if err := w.Watch(ctx); err != nil {
		logger.Fatal("watch error", zap.Error(err))
	}
}
