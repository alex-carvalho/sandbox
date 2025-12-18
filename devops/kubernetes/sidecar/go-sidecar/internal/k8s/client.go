package k8s

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	clientset     kubernetes.Interface
	namespace     string
	configMapName string
	logger        *zap.Logger
}

func NewClient(namespace, configMapName string, logger *zap.Logger) (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{
		clientset:     clientset,
		namespace:     namespace,
		configMapName: configMapName,
		logger:        logger,
	}, nil
}

// UpdateConfigMap updates or creates a ConfigMap with the given data
func (c *Client) UpdateConfigMap(ctx context.Context, data map[string]string) error {
	configMapsClient := c.clientset.CoreV1().ConfigMaps(c.namespace)

	cm, err := configMapsClient.Get(ctx, c.configMapName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      c.configMapName,
					Namespace: c.namespace,
				},
				Data: data,
			}

			_, err = configMapsClient.Create(ctx, cm, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create ConfigMap: %w", err)
			}

			c.logger.Info("created ConfigMap",
				zap.String("name", c.configMapName),
				zap.Int("entries", len(data)),
			)
			return nil
		}
		return fmt.Errorf("failed to get ConfigMap: %w", err)
	}

	cm.Data = data
	_, err = configMapsClient.Update(ctx, cm, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update ConfigMap: %w", err)
	}

	c.logger.Info("updated ConfigMap",
		zap.String("name", c.configMapName),
		zap.Int("entries", len(data)),
	)
	return nil
}

func (c *Client) GetConfigMap(ctx context.Context) (*corev1.ConfigMap, error) {
	configMapsClient := c.clientset.CoreV1().ConfigMaps(c.namespace)
	return configMapsClient.Get(ctx, c.configMapName, metav1.GetOptions{})
}
