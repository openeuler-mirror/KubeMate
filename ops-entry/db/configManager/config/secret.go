package config

import (
	"context"
	"errors"
	"fmt"
	"ops-entry/common/util"
	"ops-entry/constValue"
	"ops-entry/db/configManager"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretImpl struct {
	NameSpace  string
	SecretName string
	SecretType corev1.SecretType
	LabelData  map[string]string
}

/**
* @Description: secret初始化
* @param namespace命名空间
* @param secretName secret 名称，多版本时，需要组装最大版本信息，
* @param labelData 以label来保存多版本的secret信息
* @param secretType 默认Opaque
* return
*   @resp
*
 */

func NewSecretImpl(nameSpace, secretName string, labelData map[string]string, secretType corev1.SecretType) *SecretImpl {
	if len(nameSpace) == 0 {
		nameSpace = constValue.DefaultNameSpace
	}

	if len(secretType) < 1 {
		secretType = corev1.SecretTypeOpaque
	}

	return &SecretImpl{
		NameSpace:  nameSpace,
		SecretName: getSecretName(secretName),
		LabelData:  labelData,
		SecretType: secretType,
	}
}

func getSecretName(name string) string {
	return fmt.Sprintf("%s%s%s", constValue.Prefix, constValue.SECRET, name)
}

func (s *SecretImpl) GetListOptions(labelData map[string]string) metav1.ListOptions {
	if len(labelData) < 1 {
		return metav1.ListOptions{}
	}
	labelSelector := metav1.LabelSelector{
		MatchLabels: labelData,
	}
	return metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	}
}

func (s *SecretImpl) Get(ctx context.Context, opts metav1.GetOptions) (*corev1.Secret, error) {
	logrus.Info("s.SecretName-->", s.SecretName, s.NameSpace)
	return configManager.KCS.ClientSet.CoreV1().Secrets(s.NameSpace).Get(ctx, s.SecretName, opts)
}

func (s *SecretImpl) List(ctx context.Context, opts metav1.ListOptions) (interface{}, error) {
	return configManager.KCS.ClientSet.CoreV1().Secrets(s.NameSpace).List(ctx, opts)
}

/**
* @Description: 创建secret ，以stringData形式进行存储
* @param opts
* @param data
* return
*   @resp
*
 */

func (s *SecretImpl) Create(ctx context.Context, opts metav1.CreateOptions, data interface{}) error {
	scData, ok := data.(map[string]string)
	if !ok {
		return errors.New("secretImpl Invalid data")
	}
	if len(scData) < 1 {
		return errors.New("secretImpl Invalid empty data")
	}

	secret, err := s.Get(ctx, metav1.GetOptions{})
	if err != nil && !k8serrors.IsNotFound(err) {
		logrus.Errorf("Get secret error: %v", err)
		return err
	}

	if len(secret.Name) == 0 {
		return s.createSecret(ctx, opts, s.SecretName, scData)
	}

	newSecret := secret.DeepCopy()
	newSecret.StringData = scData
	newSecret.Type = s.SecretType
	newSecret.Labels = s.LabelData

	// 使用客户端更新Secrets
	_, err = configManager.KCS.ClientSet.CoreV1().Secrets(secret.Namespace).Update(ctx, newSecret, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("update Secrets failed: [sc:%+v],[err:%v]", newSecret, err)
		return err
	}

	return nil
}

func (s *SecretImpl) createSecret(ctx context.Context, opts metav1.CreateOptions, name string, data map[string]string) error {
	secretMap := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: s.LabelData,
		},
		StringData: data,
		Type:       s.SecretType,
	}

	_, err := configManager.KCS.ClientSet.CoreV1().Secrets(s.NameSpace).Create(ctx, secretMap, opts)
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
			logrus.Errorf("secretMap %s already exists\n", secretMap)
		} else {
			logrus.Error("secretMap create err", err)
		}
		return err
	}
	return nil
}

/**
* @Description: 添加或者修改labels
* @param cmLabels待新增或者修改的labels
* @param names存在就修改对应名称的secret的labels，不存在就修改查询结果中全部的secret
* return
*   @resp  返回最后一个修改失败的错误的信息
*    备注： 该修改不是原子性的可能存在部分修改成功，部分修改失败
 */

func (s *SecretImpl) AddLabels(ctx context.Context, opts metav1.UpdateOptions, cmLabels map[string]string, names ...string) error {
	if len(cmLabels) < 1 {
		return errors.New("secretImpl Invalid empty cmLabels")
	}

	//获取当前列表信息
	listOptions := s.GetListOptions(s.LabelData)
	cms, err := s.List(ctx, listOptions)
	if err != nil {
		return err
	}

	secretList, ok := cms.(*corev1.SecretList)
	if !ok {
		return errors.New("add  labels Secret type trans failed")
	}

	//不存在
	if len(secretList.Items) < 1 {
		return errors.New("secretList Not Exists")
	}

	for _, sr := range secretList.Items {
		if len(names) > 0 && util.IsElementInSlice(names, sr.GetName()) || len(names) == 0 {
			sr.Labels = util.MergeMap(sr.Labels, cmLabels)
			// 使用客户端更新secret
			_, err = configManager.KCS.ClientSet.CoreV1().Secrets(s.NameSpace).Update(ctx, &sr, opts)
			if err != nil {
				logrus.Errorf("update secret failed: [sr:%+v],[err:%v]", sr, err)
				continue
			}
		}
	}
	return err
}

func (s *SecretImpl) Update(ctx context.Context, opts metav1.UpdateOptions, data interface{}) error {
	scData, ok := data.(map[string]string)
	if !ok {
		return errors.New("secretImpl Update Invalid data")
	}
	if len(scData) < 1 {
		return errors.New("secretImpl Update Invalid empty data")
	}
	// 获取现有的secret对象
	listOptions := s.GetListOptions(s.LabelData)
	scs, err := s.List(ctx, listOptions)
	if err != nil {
		return err
	}

	secretList, ok := scs.(*corev1.SecretList)
	if !ok {
		return errors.New("update  secret type trans failed")
	}

	if len(secretList.Items) < 1 {
		logrus.Errorf("get empty data [label:%+v]", s.LabelData)
		return errors.New("update empty secret failed")
	}

	for _, sc := range secretList.Items {
		sc.StringData = scData
		// 使用客户端更新Secrets
		_, err = configManager.KCS.ClientSet.CoreV1().Secrets(sc.Namespace).Update(ctx, &sc, opts)
		if err != nil {
			logrus.Errorf("update Secrets failed: [sc:%+v],[err:%v]", sc, err)
			continue
		}
	}

	return err
}

func (s *SecretImpl) Delete(ctx context.Context, opts metav1.DeleteOptions) error {
	err := configManager.KCS.ClientSet.CoreV1().Secrets(s.NameSpace).Delete(ctx, s.SecretName, opts)
	if err != nil {
		logrus.Errorf("delete Secret by name failed: [name:%s],[err:%+v]", s.SecretName, err)
		return err
	}
	logrus.Infof("Secret deleted successfully: [name:%s]", s.SecretName)

	if len(s.LabelData) > 0 {
		options := s.GetListOptions(s.LabelData)
		list, err := s.List(ctx, options)
		if err != nil {
			logrus.Errorf("empty SecretImpl data:[err:%+v]", err)
			return err
		}
		secretList, ok := list.(*corev1.SecretList)
		if !ok {
			logrus.Error("type Secret conversion failed")
			return errors.New("type Secret conversion failed")
		}

		if len(secretList.Items) < 1 {
			logrus.Error("empty Secret data")
			return errors.New("empty Secret data")
		}

		// 遍历 secret 列表并执行删除操作
		for _, sc := range secretList.Items {
			err := configManager.KCS.ClientSet.CoreV1().Secrets(sc.Namespace).Delete(context.TODO(), sc.Name, opts)
			if err != nil {
				logrus.Errorf("delete Secrets failed: [cm:%+v],[err:%v]", sc, err)
				continue
			}
		}
	}

	return nil
}
