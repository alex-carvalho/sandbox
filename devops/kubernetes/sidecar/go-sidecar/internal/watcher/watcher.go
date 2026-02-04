package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"

	"github.com/example/go-sidecar/internal/k8s"
	"github.com/example/go-sidecar/internal/properties"
)

type Watcher struct {
	watchDir  string
	fsWatcher *fsnotify.Watcher
	k8sClient *k8s.Client
	logger    *zap.Logger
}

func NewWatcher(watchDir string, k8sClient *k8s.Client, logger *zap.Logger) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	w := &Watcher{
		watchDir:  watchDir,
		fsWatcher: fsWatcher,
		k8sClient: k8sClient,
		logger:    logger,
	}

	if err := fsWatcher.Add(watchDir); err != nil {
		return nil, fmt.Errorf("failed to add watch directory: %w", err)
	}

	err = filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != watchDir {
			if err := fsWatcher.Add(path); err != nil {
				w.logger.Warn("failed to watch directory", zap.String("dir", path), zap.Error(err))
			}
		}
		return nil
	})

	if err != nil {
		fsWatcher.Close()
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return w, nil
}

func (w *Watcher) Watch(ctx context.Context) error {
	if err := w.syncProperties(ctx); err != nil {
		w.logger.Error("initial sync failed", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return fmt.Errorf("watcher events channel closed")
			}

			if isPropertyFile(event.Name) {
				w.logger.Info("file event detected",
					zap.String("file", event.Name),
					zap.String("op", event.Op.String()),
				)

				if err := w.syncProperties(ctx); err != nil {
					w.logger.Error("sync failed", zap.Error(err))
				}
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return fmt.Errorf("watcher errors channel closed")
			}
			w.logger.Error("watcher error", zap.Error(err))
		}
	}
}

func (w *Watcher) syncProperties(ctx context.Context) error {
	allProps := make(map[string]string)

	err := filepath.Walk(w.watchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isPropertyFile(path) {
			file, err := os.Open(path)
			if err != nil {
				w.logger.Warn("failed to open file", zap.String("file", path), zap.Error(err))
				return nil
			}
			defer file.Close()

			props, err := properties.Parse(file)
			if err != nil {
				w.logger.Warn("failed to parse properties file", zap.String("file", path), zap.Error(err))
				return nil
			}

			for k, v := range props {
				allProps[k] = v
			}

			w.logger.Debug("loaded properties", zap.String("file", path), zap.Int("count", len(props)))
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	if len(allProps) == 0 {
		w.logger.Debug("no properties to sync")
		return nil
	}

	if err := w.k8sClient.UpdateConfigMap(ctx, allProps); err != nil {
		return fmt.Errorf("failed to update ConfigMap: %w", err)
	}

	return nil
}

func (w *Watcher) Close() error {
	return w.fsWatcher.Close()
}

func isPropertyFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".properties")
}
