package matrixhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/thomasbuchinger/homelab-api/pkg/common"
	"github.com/thomasbuchinger/homelab-api/pkg/kubernetes"
	"go.uber.org/zap"
)

type MatrixHelperConfig struct {
	ClientsecretNamespace, ClientsecretName, ClientsecretKey string
	MasNamespace, MasSecretName, MasSecretKey, MasDeployment string
	PollInterval                                             time.Duration
}

var logger *zap.SugaredLogger
var k8s *kubernetes.KubernetesClient

func NewMatrixHelperConfig() MatrixHelperConfig {
	return MatrixHelperConfig{
		ClientsecretNamespace: common.GetEnvWithDefault("MATRIX_HELPER_CLIENTSECRET_NAMESPACE", "keycloak"),
		ClientsecretName:      common.GetEnvWithDefault("MATRIX_HELPER_CLIENTSECRET_NAME", "keycloak-client-secret-client-matrix"),
		ClientsecretKey:       common.GetEnvWithDefault("MATRIX_HELPER_CLIENTSECRET_KEY", "CLIENT_SECRET"),
		MasNamespace:          "matrix",
		MasDeployment:         "matrix-matrix-authentication-service",
		MasSecretName:         "matrix-matrix-authentication-service",
		MasSecretKey:          "user-keycloak",
		PollInterval:          10 * time.Minute,
	}
}

func SyncLoop(ctx context.Context, l *zap.SugaredLogger) {
	var err error
	cfg := NewMatrixHelperConfig()
	logger = l
	k8s, err = kubernetes.NewKubernetesClient()
	if err != nil {
		panic(fmt.Errorf("creating Kubernetes client: %w", err))
	}
	for {
		syncOidcClient(cfg)
		select {
		case <-ctx.Done():
			return
		case <-time.After(cfg.PollInterval):
		}
	}
}

func syncOidcClient(cfg MatrixHelperConfig) {
	sourceSecretLabel := fmt.Sprintf("%s/%s:%s", cfg.ClientsecretNamespace, cfg.ClientsecretName, cfg.ClientsecretKey)
	masConfigLabel := fmt.Sprintf("Deployment: %s/%s | Secret: %s:%s", cfg.MasNamespace, cfg.MasDeployment, cfg.MasSecretName, cfg.MasSecretKey)

	sourceSecret, err := k8s.GetSecret(cfg.ClientsecretNamespace, cfg.ClientsecretName)
	if err != nil {
		panic(fmt.Errorf("Unable to get Source Secret: "+sourceSecretLabel+": %w", err))
	}
	masSecret, err := k8s.GetSecret(cfg.MasNamespace, cfg.MasSecretName)
	if err != nil {
		panic(fmt.Errorf("Unable to get Mas Config Secret: "+masConfigLabel+": %w", err))
	}

	desiredMasConfig, changed, err := ensureOidcClientSecret(string(masSecret.Data[cfg.MasSecretKey]), string(sourceSecret.Data[cfg.ClientsecretKey]))
	if err != nil {
		panic(fmt.Errorf("ensuring OIDC client secret | Source: %s | Mas: %s: %w", sourceSecretLabel, masConfigLabel, err))
	}

	metricCheckTimestamp.SetToCurrentTime()
	if !changed {
		logger.Infof("OK: No action needed")
		return
	}

	// Update Secret...
	masSecret.Data[cfg.MasSecretKey] = []byte(desiredMasConfig)
	if err := k8s.UpdateSecret(cfg.MasNamespace, masSecret); err != nil {
		panic(err)
	}
	metricChangeTotal.Inc()
	metricChangeTimestamp.SetToCurrentTime()

	// ... and restart MAS
	if err := k8s.RestartDeployment(cfg.MasNamespace, cfg.MasDeployment); err != nil {
		panic(err)
	}
	metricMasRestartTotal.Inc()
	metricMasRestartTimestamp.SetToCurrentTime()
	logger.Infof("OK: applied new client_secret to %s/%s and requested restart of %s/%s", cfg.MasNamespace, cfg.MasSecretName, cfg.MasNamespace, cfg.MasDeployment)
}
