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
	"errors"
	"ops-entry/constValue"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type K8sClientSet struct {
	ClientSet        *kubernetes.Clientset
	DynamicClientSet *dynamic.DynamicClient
}

var KCS *K8sClientSet

func Init() (err error) {
	KCS, err = NewK8sClientSet()
	if err != nil {
		logrus.Errorf("init K8sClientSet failed, err:%v", err)
		return err
	}

	CreateNameSpace(constValue.NameSpace)
	go KCS.CheckK8sHealth()
	return nil
}

func NewK8sClientSet() (*K8sClientSet, error) {
	kcs := &K8sClientSet{}
	cnf, err := getConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(cnf)
	if err != nil {
		logrus.Errorf("et client  failed from out configManager [err:%+v]", err.Error())
		return nil, err
	}

	// create the dynamic client
	dynamicClientSet := dynamic.NewForConfigOrDie(cnf)
	if clientSet == nil || dynamicClientSet == nil {
		logrus.Error("can't get clientSet or dynamicClientSet")
		return nil, errors.New("can't get clientSet or dynamicClientSet")
	}
	kcs.ClientSet = clientSet
	kcs.DynamicClientSet = dynamicClientSet
	return kcs, nil
}

func (kcs *K8sClientSet) CheckK8sHealth() {
	for {
		time.Sleep(500 * time.Second)
		pods, err := kcs.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logrus.Errorf("Error checking Kubernetes connectivity: %v\n", err)
			continue
		}
		logrus.Infof("Connected to Kubernetes. Found %d pods in the cluster.\n", len(pods.Items))
	}
}

func getConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" && os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		logrus.Infof("inner k8s cluster")
		// 使用集群内的配置创建clientset
		return rest.InClusterConfig()
	}
	return clientcmd.BuildConfigFromFlags("", homedir.HomeDir()+constValue.KubeConfig)
}

// GetOuterK8sConfig  获取用户上传的kubeconfig
func GetOuterK8sConfig(kubeconfigBytes []byte) (*rest.Config, error) {
	config, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	return config.ClientConfig()
}

// GetOuterClientSet  根据用户上传的kubeConfigPath获取clientSet
func GetOuterClientSet(kubeconfigBytes []byte) (*kubernetes.Clientset, error) {
	config, err := GetOuterK8sConfig(kubeconfigBytes)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// CreateNameSpace 创建命名空间
func CreateNameSpace(namespace string) {
	// 创建或检查命名空间
	_, err := KCS.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		logrus.Infof("namespace %s already exists", namespace)
		return
	}
	if !apierrors.IsNotFound(err) {
		logrus.Errorf("error creating namespace %s find err:%s", namespace, err.Error())
		return
	}

	// 命名空间不存在，创建它
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err = KCS.ClientSet.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error creating namespace %s find err:%s", namespace, err.Error())
		return
	}
	logrus.Infof("Namespace %s created.\n", namespace)
}
