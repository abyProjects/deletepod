package main

import (
	"github.com/abyProjects/deletepod/cmd"
	"flag"
	"fmt"

	// "log"
	// "strings"
	// "k8s.io/client-go/rest"

	"time"
	"log"
	"context"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	cmd.Execute()

	podName := cmd.PodName
	podNamespace := cmd.PodNamespace
	token := cmd.Token

	fmt.Println(podName, podNamespace, token)

	var deletePodName []string

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
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

	pods, err := clientset.CoreV1().Pods(podNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster under namespace: %s\n", len(pods.Items), podNamespace)
	for _, pod := range pods.Items {
		podLabel := pod.GetLabels()["app"]
		if podLabel == podName {
			deletePodName = append(deletePodName, pod.Name)
		}

	}
	fmt.Println(deletePodName)

	for{ 
		for _,pod := range(deletePodName){
			_, err = clientset.CoreV1().Pods(podNamespace).Get(context.TODO(), pod, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				fmt.Printf("Pod: %s not found in namespace: %s\n", podName, podNamespace)
			} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
				fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
			} else if err != nil {
				panic(err.Error())
			} else {
				fmt.Printf("Found pod: %s in namespace: %s\n", pod, podNamespace)

				// delete the pod
				err := clientset.CoreV1().Pods(podNamespace).Delete(context.TODO(), pod, metav1.DeleteOptions{})
				if err != nil {
					log.Fatal(err)
				}
			}

		}

		time.Sleep(10 * time.Second)
	}
	
}
