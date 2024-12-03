/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: lijian <lijian@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
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
