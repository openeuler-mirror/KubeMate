package configManager

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigInterface interface {
	GetListOptions(labels map[string]string) metav1.ListOptions
	Get(ctx context.Context, opts metav1.ListOptions) (interface{}, error)
	Create(ctx context.Context, opts metav1.CreateOptions, data interface{}) error
	Update(ctx context.Context, opts metav1.UpdateOptions, data interface{}) error
	AddLabels(ctx context.Context, opts metav1.UpdateOptions, cmLabels map[string]string, names ...string) error
	Delete(ctx context.Context, opts metav1.DeleteOptions) error
}
