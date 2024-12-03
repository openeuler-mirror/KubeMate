/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * KubeMate licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: weihuanhuan <weihuanhuan@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package config

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// 缺少gomokey mock和assert断言， gomonkey 不支持 x86/amd
func TestCrImpl(t *testing.T) {
	ctx := context.TODO()

	cr := NewCrImpl("test.io", "v1", "buyer", "buyers", "test_cr", "test_cr_name", "", nil)
	data := map[string]string{"name": "张三"}

	//mock初始化
	patches := gomonkey.ApplyFunc(cr.Get, func(ctx context.Context, opts metav1.ListOptions) (interface{}, error) {
		// 这里返回固定的数据或模拟的错误
		return &unstructured.UnstructuredList{
			Object: map[string]interface{}{},
			Items: []unstructured.Unstructured{
				{},
				{},
			},
		}, nil
	})
	defer patches.Reset() // 在测试结束后恢复原始函数

	t.Run("create", func(t *testing.T) {
		cr.LabelData = map[string]string{"test_v1": "6666"}
		err := cr.Create(ctx, metav1.CreateOptions{}, data)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
		assert.Equal(t, "failed", err)

	})

	t.Run("update", func(t *testing.T) {
		cr.LabelData = map[string]string{"test_v1": "6666"}
		data = map[string]string{"age": "44"}
		err := cr.Update(ctx, metav1.UpdateOptions{}, map[string]string{"name": "张三"})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})

	t.Run("delete", func(t *testing.T) {
		cr.LabelData = map[string]string{"test_v1": "6666"}
		data = map[string]string{"age": "44"}
		err := cr.Delete(ctx, metav1.DeleteOptions{})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
}
