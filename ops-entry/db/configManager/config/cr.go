package config

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/db/configManager"
	"strconv"
	"strings"
)

type CrImpl struct {
	Group       string
	Version     string
	Kind        string
	Resource    string
	NameSpace   string
	CrName      string
	LabelData   map[string]string
	UpdateFiled string
}

/**
* @Description: cr实例
* @param gvk
* @param gvr
* @param crname cr的名字，作为cr的唯一标识
* @param UpdateFiled 默认更新spec字段
* return
*   @resp
*
 */

func NewCrImpl(group, version, kind, resource, nameSpace, crName, updateFiled string, labelData map[string]string) *CrImpl {
	if len(nameSpace) == 0 {
		nameSpace = constValue.DefaultNameSpace
	}

	if len(version) == 0 {
		version = constValue.DefaultCrVersion
	}

	if len(updateFiled) == 0 {
		updateFiled = constValue.DefaultCrUpdateField
	}
	return &CrImpl{
		Group:       group,
		Version:     version,
		Kind:        kind,
		Resource:    resource,
		NameSpace:   nameSpace,
		CrName:      crName,
		LabelData:   labelData,
		UpdateFiled: updateFiled,
	}
}

func (c *CrImpl) GetListOptions(labels map[string]string) metav1.ListOptions {
	if len(labels) < 1 {
		return metav1.ListOptions{}
	}

	// 定义 LabelSelector 来筛选 cr
	labelSelector := metav1.LabelSelector{
		MatchLabels: labels,
	}
	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	}
}

// 获取cr的name
func (c *CrImpl) getCrName(name string, revision int) string {
	return fmt.Sprintf("%s%s%s_%s%d", constValue.Prefix, constValue.Cr, name, constValue.VersionMark, revision)
}

func (c *CrImpl) Get(ctx context.Context, opts metav1.ListOptions) (interface{}, error) {
	gvr := schema.GroupVersionResource{
		Group:    c.Group,
		Version:  c.Version,
		Resource: c.Resource,
	}

	// 获取自定义资源列表
	return configManager.KCS.DynamicClientSet.Resource(gvr).Namespace(c.NameSpace).List(ctx, opts)
}

/**
* @Description: cr实例创建
* @param opts CreateOptions
* @param data需要设置的数据
* @param 修改的默认列为spec
* return
*   @resp
*
 */

func (c *CrImpl) Create(ctx context.Context, opts metav1.CreateOptions, data interface{}) error {
	crData, ok := data.(map[string]string)
	if !ok {
		return errors.New("CrImpl Invalid data")
	}
	if len(crData) < 1 {
		return errors.New("CrImpl Invalid empty data")
	}

	listOptions := c.GetListOptions(c.LabelData)

	gvr := schema.GroupVersionResource{
		Group:    c.Group,
		Version:  c.Version,
		Resource: c.Resource,
	}
	gvk := schema.GroupVersionKind{
		Group:   c.Group,
		Version: c.Version,
		Kind:    c.Kind,
	}
	// 获取自定义资源列表
	list, err := c.Get(ctx, listOptions)
	if err != nil {
		return err
	}

	crList, ok := list.(*unstructured.UnstructuredList)
	if !ok {
		return errors.New("Create  cr type trans failed")
	}
	if len(crList.Items) < 1 {
		name := c.getCrName(c.CrName, 1)
		return c.createCr(ctx, gvk, gvr, opts, name, crData)
	}

	// 获取当前最新revision Data字段信息
	index := 0
	for _, cr := range crList.Items {
		idx := strings.LastIndex(cr.GetName(), constValue.VersionMark)
		if idx < 0 {
			logrus.Errorf("create [cm:%+v]", cr)
			return errors.New("mapImpl Invalid data " + cr.GetName())
		}
		revision, err := strconv.Atoi(cr.GetName()[idx+len(constValue.VersionMark):])
		if err != nil {
			return err
		}
		if revision > index {
			index = revision
		}
	}
	name := c.getCrName(c.CrName, index+1)
	return c.createCr(ctx, gvk, gvr, opts, name, crData)
}

func (c *CrImpl) createCr(ctx context.Context, gvk schema.GroupVersionKind, gvr schema.GroupVersionResource, opts metav1.CreateOptions, name string, crData map[string]string) error {
	// 创建CR实例
	cr := &unstructured.Unstructured{}
	cr.SetGroupVersionKind(gvk)
	cr.SetNamespace(c.NameSpace) // 设置命名空间
	cr.SetName(name)             // 设置CR名称
	cr.SetLabels(c.LabelData)    // 设置labels
	// 设置CR的Spec字段
	for k, val := range crData {
		err := unstructured.SetNestedField(cr.Object, val, c.UpdateFiled, k)
		if err != nil {
			logrus.Errorf("cr set value failed [err:%v],[data:%+v]", err, crData)
			return err
		}
	}

	// 创建CR
	_, err := configManager.KCS.DynamicClientSet.Resource(gvr).Namespace(c.NameSpace).Create(ctx, cr, opts)
	if err != nil {
		logrus.Errorf("cr create failed [err:%v]", err)
		return err
	}
	return nil
}

/**
* @Description: 添加或者修改labels
* @param cmLabels待新增或者修改的labels
* @param names存在就修改对应名称的cr的labels，不存在就修改查询结果中全部的cr
* return
*   @resp  返回最后一个修改失败的错误的信息
*    备注： 该修改不是原子性的可能存在部分修改成功，部分修改失败
 */

func (c *CrImpl) AddLabels(ctx context.Context, opts metav1.UpdateOptions, cmLabels map[string]string, names ...string) error {
	if len(cmLabels) < 1 {
		return errors.New("crImpl Invalid empty cmLabels")
	}

	//获取当前列表信息
	listOptions := c.GetListOptions(c.LabelData)

	gvr := schema.GroupVersionResource{
		Group:    c.Group,
		Version:  c.Version,
		Resource: c.Resource,
	}

	// 获取自定义资源列表
	list, err := c.Get(ctx, listOptions)
	if err != nil {
		return err
	}

	listItem, ok := list.(*unstructured.UnstructuredList)
	if !ok {
		return errors.New("invalid data" + c.CrName)
	}

	if len(listItem.Items) < 1 {
		logrus.Errorf("empty cr data [label:%+v]", c.LabelData)
		return errors.New("empty cr data")
	}

	for _, item := range listItem.Items {
		if len(names) > 0 && util.IsElementInSlice(names, item.GetName()) || len(names) == 0 {
			labels := util.MergeMap(item.GetLabels(), cmLabels)
			item.SetLabels(labels)
			_, err = configManager.KCS.DynamicClientSet.Resource(gvr).Namespace(c.NameSpace).Update(ctx, &item, opts)
			return err
		}
	}
	return err
}

func (c *CrImpl) Update(ctx context.Context, opts metav1.UpdateOptions, data interface{}) error {
	crData, ok := data.(map[string]string)
	if !ok {
		return errors.New("CrImpl Invalid data")
	}
	if len(crData) < 1 {
		return errors.New("CrImpl Invalid empty data")
	}

	listOptions := c.GetListOptions(c.LabelData)

	gvr := schema.GroupVersionResource{
		Group:    c.Group,
		Version:  c.Version,
		Resource: c.Resource,
	}

	// 获取自定义资源列表
	list, err := c.Get(ctx, listOptions)
	if err != nil {
		return err
	}

	listItem, ok := list.(*unstructured.UnstructuredList)
	if !ok {
		return errors.New("invalid data" + c.CrName)
	}

	if len(listItem.Items) < 1 {
		logrus.Errorf("empty cr data [label:%+v]", c.LabelData)
		return errors.New("empty cr data")
	}

	for _, item := range listItem.Items {
		for k, val := range crData {
			err = unstructured.SetNestedField(item.Object, val, c.UpdateFiled, k)
			if err != nil {
				logrus.Errorf("cr init value failed [err:%v],[data:%+v]", err, crData)
				return err
			}
		}
		_, err = configManager.KCS.DynamicClientSet.Resource(gvr).Namespace(c.NameSpace).Update(ctx, &item, opts)
		return err

	}
	return nil
}

func (c *CrImpl) Delete(ctx context.Context, opts metav1.DeleteOptions) error {
	gvr := schema.GroupVersionResource{
		Group:    c.Group,
		Version:  c.Version,
		Resource: c.Resource,
	}
	options := c.GetListOptions(c.LabelData)
	list, err := c.Get(ctx, options)
	if err != nil {
		logrus.Errorf("empty cr data：[err:%+v]", err)
		return err
	}
	listItem, ok := list.(*unstructured.UnstructuredList)
	if !ok {
		logrus.Errorf("empty cr data：[labels:%+v]", c.LabelData)
		return errors.New("invalid data")
	}

	if len(listItem.Items) < 1 {
		logrus.Errorf("empty cr data [label:%+v]", c.LabelData)
		return errors.New("empty cr data")
	}

	// 遍历 cr 列表并执行删除操作
	for _, cr := range listItem.Items {
		err := configManager.KCS.DynamicClientSet.Resource(gvr).Namespace(c.NameSpace).Delete(context.TODO(), cr.GetName(), opts)
		if err != nil {
			logrus.Errorf("delete cr failed: [cr:%+v],[err:%v]", cr, err)
			continue
		}
	}
	return err
}
