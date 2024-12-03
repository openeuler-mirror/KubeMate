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
	"fmt"
	"ops-entry/constValue"
	"strconv"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"errors"
	"ops-entry/common/util"
	"ops-entry/db/configManager"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MapImpl struct {
	NameSpace     string
	ConfigMapName string
	LabelData     map[string]string
}

/**
* @Description:
* @param nameSpace命名空间，默认为default
* @param 自定义的configMap的名称
* @param labelData label数据列表必须要有一个标识来进行多版本管理
* return
*   @resp configMap实例
*
 */

func NewMapImpl(nameSpace, configMapName string, labelData map[string]string) *MapImpl {
	if len(nameSpace) == 0 {
		nameSpace = constValue.DefaultNameSpace
	}
	return &MapImpl{
		NameSpace:     nameSpace,
		ConfigMapName: configMapName,
		LabelData:     labelData,
	}
}

/**
* @Description:
* @param labelData  label数据
* return
*   @resp opts
*
 */

func (m *MapImpl) GetListOptions(labelData map[string]string) metav1.ListOptions {
	if len(labelData) < 1 {
		return metav1.ListOptions{}
	}

	// 定义 LabelSelector 来筛选 ConfigMap
	labelSelector := metav1.LabelSelector{
		MatchLabels: labelData,
	}
	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	}
}

// 获取configmap的name
func (m *MapImpl) getConfigMapName(name string, revision int) string {
	return fmt.Sprintf("%s%s%s_%s%d", constValue.Prefix, constValue.ConfigMap, name, constValue.VersionMark, revision)
}

/**
* @Description:
* @param opts labels列表
* return
*   @resp  返回configMap 列表
*
 */

func (m *MapImpl) Get(ctx context.Context, opts metav1.ListOptions) (interface{}, error) {
	fmt.Println(99999)
	// 列出具有特定标签的 ConfigMap
	return configManager.KCS.ClientSet.CoreV1().ConfigMaps(m.NameSpace).List(ctx, opts)
}

/**
* @Description: 添加或者修改labels
* @param cmLabels待新增或者修改的labels
* @param names存在就修改对应名称的configMap的labels，不存在就修改查询结果中全部的configMap
* return
*   @resp  返回最后一个修改失败的错误的信息
*    备注： 该修改不是原子性的可能存在部分修改成功，部分修改失败
 */

func (m *MapImpl) AddLabels(ctx context.Context, opts metav1.UpdateOptions, cmLabels map[string]string, names ...string) error {
	if len(cmLabels) < 1 {
		return errors.New("mapImpl Invalid empty cmLabels")
	}

	//获取当前列表信息
	listOptions := m.GetListOptions(m.LabelData)
	cms, err := m.Get(ctx, listOptions)
	if err != nil {
		return err
	}

	configMapList, ok := cms.(*corev1.ConfigMapList)
	if !ok {
		return errors.New("add  labels configMap type trans failed")
	}

	//不存在
	if len(configMapList.Items) < 1 {
		return errors.New(" configMap Not Exists")
	}

	for _, cm := range configMapList.Items {
		if len(names) > 0 && util.IsElementInSlice(names, cm.GetName()) || len(names) == 0 {
			cm.Labels = util.MergeMap(cm.Labels, cmLabels)
			// 使用客户端更新ConfigMap
			_, err = configManager.KCS.ClientSet.CoreV1().ConfigMaps(m.NameSpace).Update(ctx, &cm, opts)
			if err != nil {
				logrus.Errorf("update configMap failed: [cm:%+v],[err:%v]", cm, err)
				continue
			}
		}
	}
	return err
}

/**
* @Description:configmap 的创建，只能保存非加密信息,需要进行多版本管理
* name_v1,name_v2 ... name_v10
* @param opts CreateOptions列表
* @param data 创建configMap时的map列表，可以为普通的key:value,也可为  文件名称：文件内容
* return
*   @resp  返回configMap创建失败的具体信息
*
 */

func (m *MapImpl) Create(ctx context.Context, opts metav1.CreateOptions, data interface{}) error {
	cmData, ok := data.(map[string]string)
	if !ok {
		return errors.New("mapImpl Invalid data")
	}
	if len(cmData) < 1 {
		return errors.New("mapImpl Invalid empty data")
	}

	//获取当前列表信息
	listOptions := m.GetListOptions(m.LabelData)
	cms, err := m.Get(ctx, listOptions)
	if err != nil {
		return err
	}

	configMapList, ok := cms.(*corev1.ConfigMapList)
	if !ok {
		return errors.New("Create  configMap type trans failed")
	}

	//初次创建
	if len(configMapList.Items) < 1 {
		name := m.getConfigMapName(m.ConfigMapName, 1)
		return m.createConfigMap(ctx, opts, name, cmData)
	}

	// 获取当前最新revision Data字段信息
	index := 0
	for _, cm := range configMapList.Items {
		idx := strings.LastIndex(cm.Name, constValue.VersionMark)
		if idx < 0 {
			logrus.Errorf("create [cm:%+v]", cm)
			return errors.New("mapImpl Invalid data " + cm.Name)
		}
		revision, err := strconv.Atoi(cm.Name[idx+len(constValue.VersionMark):])
		if err != nil {
			return err
		}
		if revision > index {
			index = revision
		}
	}

	name := m.getConfigMapName(m.ConfigMapName, index+1)
	return m.createConfigMap(ctx, opts, name, cmData)
}

/**
* @Description: 创建configMap 多版本区分标识为，configMap的名称不一样，但是label一样
* @param name configMap的名称
* @param data待保存的数据
* @param
* @param
* return
*   @resp
*
 */

func (m *MapImpl) createConfigMap(ctx context.Context, opts metav1.CreateOptions, name string, data map[string]string) error {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: m.LabelData,
		},
		Data: data,
	}

	_, err := configManager.KCS.ClientSet.CoreV1().ConfigMaps(m.NameSpace).Create(ctx, configMap, opts)
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			logrus.Errorf("ConfigMap %s already exists\n", configMap)
		} else {
			logrus.Error("ConfigMap create err", err)
		}
		return err
	}
	return nil
}

/**
* @Description: 更新configMap列表,存在多版本时，需要全部进行修改
* @param ctx
* @param opts
* return
*   @resp 可能更新多个configMap,不能中断，只返回最后一个错误；日志文件中包含全部错误信息
*
 */

func (m *MapImpl) Update(ctx context.Context, opts metav1.UpdateOptions, data interface{}) error {
	cmData, ok := data.(map[string]string)
	if !ok {
		return errors.New("mapImpl Update Invalid data")
	}
	if len(cmData) < 1 {
		return errors.New("mapImpl Update Invalid empty data")
	}
	// 获取现有的ConfigMap对象
	listOptions := m.GetListOptions(m.LabelData)
	cms, err := m.Get(ctx, listOptions)
	if err != nil {
		return err
	}

	configMapList, ok := cms.(*corev1.ConfigMapList)
	if !ok {
		return errors.New("update  configMap type trans failed")
	}

	if len(configMapList.Items) < 1 {
		logrus.Errorf("get empty data [label:%+v]", m.LabelData)
		return errors.New("update empty configMap failed")
	}

	for _, cm := range configMapList.Items {
		cm.Data = util.MergeMap(cm.Data, cmData)
		// 使用客户端更新ConfigMap
		_, err = configManager.KCS.ClientSet.CoreV1().ConfigMaps(m.NameSpace).Update(ctx, &cm, opts)
		if err != nil {
			logrus.Errorf("update configMap failed: [cm:%+v],[err:%v]", cm, err)
			continue
		}
	}

	return err
}

/**
* @Description: 删除configMap列表
* @param ctx
* @param opts
* return
*   @resp 可能删除多个configMap,不能中断，只返回最后一个错误；日志文件中包含全部错误信息
*
 */

func (m *MapImpl) Delete(ctx context.Context, opts metav1.DeleteOptions) error {
	options := m.GetListOptions(m.LabelData)
	list, err := m.Get(ctx, options)
	if err != nil {
		logrus.Errorf("empty configMaps data：[err:%+v]", err)
		return err
	}
	configMapList, ok := list.(*corev1.ConfigMapList)
	if !ok {
		logrus.Error("type ConfigMap conversion failed")
		return errors.New("type ConfigMap conversion failed")
	}

	if len(configMapList.Items) < 1 {
		logrus.Error("empty configMaps data")
		return errors.New("empty configMaps data")
	}

	// 遍历 ConfigMap 列表并执行删除操作
	for _, cm := range configMapList.Items {
		err := configManager.KCS.ClientSet.CoreV1().ConfigMaps(cm.Namespace).Delete(context.TODO(), cm.Name, opts)
		if err != nil {
			logrus.Errorf("delete configMap failed: [cm:%+v],[err:%v]", cm, err)
			continue
		}
	}
	return err
}
