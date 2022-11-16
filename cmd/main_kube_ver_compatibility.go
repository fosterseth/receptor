package main

import (
  "fmt"
	"k8s.io/client-go/discovery"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/tools/clientcmd"
)
 
func main() {
  clr := clientcmd.NewDefaultClientConfigLoadingRules()
  config, err := clientcmd.BuildConfigFromFlags("", clr.GetDefaultFilename())
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return
	}
	serverVerInfo, err := dc.ServerVersion()
	if err != nil {
		return
	}
	fmt.Printf("Kubernetes Server Version: %#v\n", serverVerInfo)
	if version.CompareKubeAwareVersionStrings("v1.25.4", serverVerInfo.GitVersion) > 0 {
    fmt.Println("will use timestamps method")
		return
	}
	var compatibleVer string
	if serverVerInfo.Minor == "23" {
		compatibleVer = "v1.23.14"
	} else if serverVerInfo.Minor == "24" {
		compatibleVer = "v1.24.8"
	}
	verDiff := version.CompareKubeAwareVersionStrings(compatibleVer, serverVerInfo.GitVersion)
	if verDiff < 0 {
		fmt.Printf("Kubernetes version %s not at least %s, using legacy io.Copy\n", serverVerInfo.GitVersion, compatibleVer)
		
		return
	}
  fmt.Println("will use timestamps method")
}
