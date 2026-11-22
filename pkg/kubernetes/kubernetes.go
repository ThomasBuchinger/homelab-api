package kubernetes

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type KubernetesClient struct {
	k8sClient  *kubernetes.Clientset
	restConfig *rest.Config
	Logger     *zap.SugaredLogger
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
		k8sClient:  clientset,
		restConfig: config,
		Logger:     common.GetServerConfig().RootLogger.Named("Kubernetes"),
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

func (k *KubernetesClient) GetSecret(namespace, name string) (*corev1.Secret, error) {
	secret, err := k.k8sClient.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		k.Logger.Errorf("Failed to Get Secret: %s/%s: %s", namespace, name, err.Error())
		return nil, err
	}
	return secret, nil
}

func (k *KubernetesClient) UpdateSecret(namespace string, secret *corev1.Secret) error {
	_, err := k.k8sClient.CoreV1().Secrets(namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		k.Logger.Errorf("Failed to Update Secret: %s/%s: %s", namespace, secret.Name, err.Error())
		return err
	}

	k.Logger.Infof("OK: Updated Secret: %s/%s", namespace, secret.Name)
	return nil
}

func (k *KubernetesClient) GetPodByLabel(namespace, labelSelector string) (podNamespace, podName string, err error) {
	pods, err := k.k8sClient.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		k.Logger.Errorf("Failed to List Pods: namespace=%s selector=%q: %s", namespace, labelSelector, err.Error())
		return "", "", err
	}
	if len(pods.Items) == 0 {
		k.Logger.Errorf("No Pods found: namespace=%s selector=%q", namespace, labelSelector)
		return "", "", fmt.Errorf("no pods found matching %q in namespace %s", labelSelector, namespace)
	}
	return pods.Items[0].Namespace, pods.Items[0].Name, nil
}

func (k *KubernetesClient) ExecInPod(podNamespace, podName, container string, command []string) (stdout, stderr string, err error) {
	req := k.k8sClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNamespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdout:    true,
			Stderr:    true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		k.Logger.Errorf("Failed to create exec executor: pod=%s/%s container=%s: %s", podNamespace, podName, container, err.Error())
		return "", "", err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	err = executor.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
	})
	if err != nil {
		k.Logger.Errorf("Failed to exec in Pod: pod=%s/%s container=%s command=%v: %s", podNamespace, podName, container, command, err.Error())
		return stdoutBuf.String(), stderrBuf.String(), err
	}
	return stdoutBuf.String(), stderrBuf.String(), nil
}
