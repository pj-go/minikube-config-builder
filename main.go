package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

var (
	kubeconfig            string
	apiServer             string
	DefaultKubeconfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Paths to a kubeconfig.")
	flag.StringVar(&apiServer, "apiServer", "", "format: https://ip:port")
}

func main() {
	flag.Parse()
	if len(kubeconfig) == 0 {
		kubeconfig = DefaultKubeconfigPath
	}
	bs, err := os.ReadFile(kubeconfig)
	if err != nil {
		panic(err)
	}

	config := &v1.Config{}
	err = yaml.Unmarshal(bs, config)
	if err != nil {
		panic(err)
	}

	spew.Dump(config)
	fmt.Println("------")

	for i, cluster := range config.Clusters {
		if cluster.Name == "minikube" {
			cluster.Cluster.CertificateAuthority = ""
			cluster.Cluster.InsecureSkipTLSVerify = true
			if len(apiServer) > 0 {
				cluster.Cluster.Server = apiServer
			}
		}
		config.Clusters[i] = cluster
	}

	for i, info := range config.AuthInfos {
		if info.Name == "minikube" {
			if len(info.AuthInfo.ClientCertificate) > 0 {
				s, err := EncodeDataFromFile(info.AuthInfo.ClientCertificate)
				if err != nil {
					panic(s)
				}
				info.AuthInfo.ClientCertificateData = []byte(s)
				info.AuthInfo.ClientCertificate = ""
			}
			if len(info.AuthInfo.ClientKey) > 0 {
				s, err := EncodeDataFromFile(info.AuthInfo.ClientKey)
				if err != nil {
					panic(s)
				}
				info.AuthInfo.ClientKeyData = []byte(s)
				info.AuthInfo.ClientKey = ""
			}
		}
		config.AuthInfos[i] = info
	}

	rbs, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(kubeconfig, rbs, 0600)
	if err != nil {
		panic(err)
	}
}

func EncodeDataFromFile(p string) (string, error) {
	bs, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(bs), err
}
