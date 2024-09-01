package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesClient struct {
	k8sClient *kubernetes.Clientset
	Logger    *zap.SugaredLogger
}

func NewKubernetesClient() (*KubernetesClient, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &KubernetesClient{
		k8sClient: clientset,
		Logger:    common.GetServerConfig().RootLogger.Named("Kubernetes"),
	}, nil
}

func (k *KubernetesClient) RestartDeployment(namespace, name string) error {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().String())
	_, err := k.k8sClient.AppsV1().Deployments(namespace).Patch(context.Background(), name, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		k.Logger.Errorf("Failed to Restart Deployment: %s/%s: %s", namespace, name, err.Error())
		return err
	}

	k.Logger.Infof("OK: Requested Restart of %s/%s", namespace, name)
	return nil
}
