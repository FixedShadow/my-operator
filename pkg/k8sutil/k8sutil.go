package k8sutil

import (
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/url"
	"os"
)

const KubeConfigEnv = "KUBECONFIG"

type ClusterConfig struct {
	Host           string
	TLSConfig      rest.TLSClientConfig
	AsUser         string
	KubeconfigPath string
}

func NewClusterConfig(config ClusterConfig) (*rest.Config, error) {
	var cfg *rest.Config
	var err error
	if config.KubeconfigPath == "" {
		config.KubeconfigPath = os.Getenv(KubeConfigEnv)
	}

	if config.KubeconfigPath != "" {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		cfg, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, fmt.Errorf("error creating config from %s: %w", config.KubeconfigPath, err)
		}
	} else {
		if len(config.Host) == 0 {
			if cfg, err = rest.InClusterConfig(); err != nil {
				return nil, err
			}
		} else {
			cfg = &rest.Config{
				Host: config.Host,
			}
			hostURL, err := url.Parse(config.Host)
			if err != nil {
				return nil, fmt.Errorf("error parsing host url %s: %w", config.Host, err)
			}
			if hostURL.Scheme == "https" {
				cfg.TLSClientConfig = config.TLSConfig
			}
		}
	}
	cfg.QPS = 100
	cfg.Burst = 100
	cfg.UserAgent = fmt.Sprintf("my-operator-test")
	cfg.Impersonate.UserName = config.AsUser
	return cfg, nil
}
