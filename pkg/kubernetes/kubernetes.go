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

	k.Logger.Infof("KO: Requested Restart of %s/%s", namespace, name)
	return nil
}

// func main() {
// 	// creates the in-cluster config
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// creates the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	test := v1.Deployment("as", "bbb")
// 	deployment, _ := clientset.AppsV1().Deployments("").Get(context.TODO(), "syncthing", metav1.GetOptions{})

// 	for {
// 		// get pods in all the namespaces by omitting namespace
// 		// Or specify namespace to get pods in particular namespace
// 		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

// 		// Examples for error handling:
// 		// - Use helper functions e.g. errors.IsNotFound()
// 		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
// 		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
// 		if errors.IsNotFound(err) {
// 			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
// 		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
// 			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
// 		} else if err != nil {
// 			panic(err.Error())
// 		} else {
// 			fmt.Printf("Found example-xxxxx pod in default namespace\n")
// 		}

// 		time.Sleep(10 * time.Second)
// 	}
// }
