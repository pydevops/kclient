/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func main() {
	// boilerplate to get clientset object
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// list nodes
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, node := range nodes.Items {
		fmt.Printf("Node: %s:%s:%v\n", node.Name, node.Spec.PodCIDR, node.Status.Addresses)
		//fmt.Printf("NodeInfo: %s\n", node.Status.NodeInfo)
	}
	namespace := "default"
	// get pods
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	for _, pod := range pods.Items {
		if pod.Status.Phase != "Running" {
			fmt.Printf("pod: %s %v %s %s\n", pod.Name, pod.Spec.NodeName, pod.Status.Phase, pod.Status.PodIP)
		}
	}

	services, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, service := range services.Items {
		fmt.Printf("Service: %s, ClusterIP:%s, ExternalIPs:%v \n", service.Name, service.Spec.ClusterIP, service.Spec.ExternalIPs)
		for _, port := range service.Spec.Ports {
			fmt.Printf("\t Port - %d\n", port.Port)
		}
	}
	//time.Sleep(10 * time.Second)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
