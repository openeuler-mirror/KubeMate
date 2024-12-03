/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: liukuo <liukuo@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package config

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// 缺少gomokey mock和assert断言
func TestSecretImpl(t *testing.T) {
	ctx := context.TODO()

	sr := NewSecretImpl("test_np_secret", "test_secret", nil, corev1.SecretTypeOpaque)
	data := map[string]string{"user": "张三", "pwd": "123456"}

	//mock初始化 go test -gcflags=all=-l
	patches := gomonkey.ApplyFunc(sr.Get, func(ctx context.Context, opts metav1.ListOptions) (interface{}, error) {
		// 这里返回固定的数据或模拟的错误
		return &corev1.SecretList{
			Items: []corev1.Secret{
				{},
				{},
			},
		}, nil
	})
	defer patches.Reset() // 在测试结束后恢复原始函数

	t.Run("create", func(t *testing.T) {
		sr.LabelData = map[string]string{"test_v1": "6666"}
		err := sr.Create(ctx, metav1.CreateOptions{}, data)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})

	t.Run("update", func(t *testing.T) {
		sr.LabelData = map[string]string{"test_v1": "6666"}
		data = map[string]string{"age": "44"}
		err := sr.Update(ctx, metav1.UpdateOptions{}, map[string]string{"name": "张三"})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})

	t.Run("delete", func(t *testing.T) {
		sr.LabelData = map[string]string{"test_v1": "6666"}
		data = map[string]string{"age": "44"}
		err := sr.Delete(ctx, metav1.DeleteOptions{})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
}
