package main

import (
	"flag"
	"os"
	"io/ioutil"

	"github.com/abyProjects/deletepod/cmd"

	"context"
	"log"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const expectedToken = "abcd"

func inOrOutClusterCheck() bool {
	var kubeconfig string

	if home,_ := os.UserHomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return false
	}

	_, err := os.Open(kubeconfig)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func main() {
	cmd.Execute()

	podName := cmd.PodName
	podNamespace := cmd.PodNamespace
	token := cmd.Token

	// check if token provided is file path
	if _, err := os.Stat(token); err == nil {
		fileContent, err := ioutil.ReadFile(token)
		if err != nil {
			log.Println("couldn't read from file")
			log.Fatal(err)
		}
		fileContentString := string(fileContent)
		if fileContentString != expectedToken {
			log.Fatalln("user authentication failed")
		}
	}else {
		if token != expectedToken {
			log.Fatalln("user authentication failed")
		}
	}
	log.Println("user authentication successfull")

	var (
		deletePodName []string
		kubeconfig    *string
		config        *rest.Config
		clientset     *kubernetes.Clientset
		err           error
	)

	if inOrOutClusterCheck() {
		log.Println("out cluster")
		kubeconfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		flag.Parse()

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}

	} else {
		log.Println("in cluster")
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	for {
		pods, err := clientset.CoreV1().Pods(podNamespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		log.Printf("found %d pod(s) in the cluster under namespace: %s\n", len(pods.Items), podNamespace)

		count := 0
		for _, pod := range pods.Items {
			if len(pods.Items) < 1 {
				break
			} else {
				podLabel := pod.GetLabels()["app"]
				if podLabel == podName {
					deletePodName = append(deletePodName, pod.Name)
					count += 1
				}
			}
		}
		if count > 0 {
			log.Printf("found %d pod(s) with name %s", count, podName)
			log.Println(deletePodName)
		} else {
			log.Printf("%d pod(s) with name %s", count, podName)
		}

		for _, pod := range deletePodName {
			if len(deletePodName) < 1 {
				log.Printf("Pod: %s not found in namespace: %s\n", podName, podNamespace)
				break
			} else {
				_, err = clientset.CoreV1().Pods(podNamespace).Get(context.TODO(), pod, metav1.GetOptions{})
				if errors.IsNotFound(err) {
					log.Printf("Pod: %s not found in namespace: %s\n", podName, podNamespace)
				} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
					log.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
				} else if err != nil {
					panic(err.Error())
				} else {
					log.Printf("Found pod: %s in namespace: %s\n", pod, podNamespace)

					// delete the pod
					err := clientset.CoreV1().Pods(podNamespace).Delete(context.TODO(), pod, metav1.DeleteOptions{})
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("pod: %s is deleted from namespace: %s\n", pod, podNamespace)
					deletePodName = []string{}
				}
			}

		}

		time.Sleep(10 * time.Second)
	}

}
