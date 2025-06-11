package main

import (
	"context"
	"flag"
	"github.com/FixedShadow/my-operator/pkg/k8sutil"
	"github.com/FixedShadow/my-operator/pkg/server"
	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	infoLogger      = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	errorLogger     = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	tlsClientConfig rest.TLSClientConfig
	impersonateUser string
	apiServer       string
)

func parseFlags(fs *flag.FlagSet) {
	_ = fs.Parse(os.Args[1:])
}

func run(fs *flag.FlagSet) int {
	parseFlags(fs)
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)
	restConfig, err := k8sutil.NewClusterConfig(k8sutil.ClusterConfig{
		Host:      apiServer,
		TLSConfig: tlsClientConfig,
		AsUser:    impersonateUser,
	})
	if err != nil {
		errorLogger.Println("failed to create Kubernetes client configuration ", err.Error())
		cancel()
		return 1
	}
	kclient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		errorLogger.Println("failed to create Kubernetes client", err.Error())
		cancel()
		return 1

	}
	kubernetesVersion, err := kclient.Discovery().ServerVersion()
	if err != nil {
		errorLogger.Println("failed to request Kubernetes server version ", err.Error())
		cancel()
		return 1
	}
	infoLogger.Println("kubernetes version: ", kubernetesVersion.String())
	mux := http.NewServeMux()
	mux.Handle("/ping", http.HandlerFunc(ping))
	srv, err := server.NewServer(mux)
	if err != nil {
		errorLogger.Println("failed to create web server ", err.Error())
		cancel()
		return 1
	}
	wg.Go(func() error { return srv.Serve() })
	infoLogger.Println("start http server successfully.")
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		infoLogger.Println("received SIGTERM, exiting gracefully...")
	case <-ctx.Done():
	}

	cancel()
	if err := wg.Wait(); err != nil {
		errorLogger.Println("unhandled error received, Exiting... ", err.Error())
		return 1
	}
	return 0
}

func ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	os.Exit(run(flag.CommandLine))
}
