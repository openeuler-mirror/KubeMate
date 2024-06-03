package config

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/db"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// 缺少gomokey mock和assert断言
func TestMapImpl(t *testing.T) {

	db.InitDb()
	ctx := context.TODO()

	configMap := NewMapImpl("test_np_config_map", "test_config_map", nil)
	data := map[string]string{"name": "张三"}

	//mock初始化
	patches := gomonkey.ApplyMethod(reflect.TypeOf(configMap), "Get", func() (interface{}, error) {
		fmt.Println("patches doing")
		// 这里返回固定的数据或模拟的错误
		return &corev1.ConfigMapList{
			Items: []corev1.ConfigMap{
				{
					Data: map[string]string{"name": "lucy", "from": "china"},
				},
				{
					Data: map[string]string{"name": "lily", "from": "china"},
				},
			},
		}, nil
	})
	defer patches.Reset() // 在测试结束后恢复原始函数

	t.Run("create", func(t *testing.T) {
		cms, err2 := configMap.Get(ctx, metav1.ListOptions{})
		if err2 != nil {
			t.Error("get err", err2)
			return
		}

		configMapList, ok := cms.(*corev1.ConfigMapList)

		if !ok {
			t.Error("type trans ", err2)
			return
		}
		configMapList.Items = nil
		//初次创建
		if len(configMapList.Items) < 1 {
			name := configMap.getConfigMapName(configMap.ConfigMapName, 1)
			t.Log("name-->", name)
			configMap.createConfigMap(ctx, metav1.CreateOptions{}, name, data)
		}
		configMapList.Items = []corev1.ConfigMap{
			{
				Data: map[string]string{"name": "lucy"},
				ObjectMeta: metav1.ObjectMeta{
					Name: "KCOD_CONFIG_MAP_test_config_map_V_4",
				},
			},
			{
				Data: map[string]string{"name": "lily"},
				ObjectMeta: metav1.ObjectMeta{
					Name: "KCOD_CONFIG_MAP_test_config_map_V_3",
				},
			},
		}

		index := 0
		for _, cm := range configMapList.Items {
			idx := strings.LastIndex(cm.Name, constValue.VersionMark)
			if idx < 0 {
				t.Errorf("create [cm:%+v]", cm)
				return
			}
			revision, err := strconv.Atoi(cm.Name[idx+len(constValue.VersionMark):])
			if err != nil {
				t.Errorf("create [cm:%+v]", err)
				return
			}
			if revision > index {
				index = revision
			}
		}

		name := configMap.getConfigMapName(configMap.ConfigMapName, index+1)
		t.Log("name-->", name)
		assert.Equal(t, "KCOD_CONFIG_MAP_test_config_map_V_5", name)

	})

	t.Run("update", func(t *testing.T) {
		configMap.LabelData = map[string]string{"test_v1": "6666"}
		data = map[string]string{"age": "44"}
		cms, err := configMap.Get(ctx, metav1.ListOptions{})
		if err != nil {
			t.Error(err)
			return
		}

		configMapList, ok := cms.(*corev1.ConfigMapList)
		if !ok {
			t.Error("update  configMap type trans failed")
			return
		}

		if len(configMapList.Items) < 1 {
			t.Errorf("get empty data [label:%+v]", nil)
			return
		}

		data = map[string]string{"name": "张三", "age": "44"}

		for _, cm := range configMapList.Items {
			cm.Data = util.MergeMap(cm.Data, data)
			t.Log(cm.Data)
			assert.Equal(t, data, cm.Data)
			// 使用客户端更新ConfigMap
			//_, err = configManager.KCS.ClientSet.CoreV1().ConfigMaps(m.NameSpace).Update(ctx, &cm, opts)
			//if err != nil {
			//	logrus.Errorf("update configMap failed: [cm:%+v],[err:%v]", cm, err)
			//	continue
			//}
		}
	})

	t.Run("delete", func(t *testing.T) {
		configMap.LabelData = map[string]string{"test_v1": "6666"}
		err := configMap.Delete(ctx, metav1.DeleteOptions{})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("success")
	})
}
