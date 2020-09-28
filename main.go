package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset
var podName string
var podPort string
var namespace string

func FetchIPsFromCluster(podName, namespace string) []string {
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ips := []string{}

	if len(pods.Items) >= 1 {
		for _, pod := range pods.Items {
			if len(pod.Name) >= len(podName) {
				if pod.Name[:len(podName)] == podName {
					ips = append(ips, pod.Status.HostIP)
				}
			}
		}
	}

	return ips
}

func exporterHandler(w http.ResponseWriter, r *http.Request) {
	resp := ""
	var auxiliar string
	var array []string

	ips := FetchIPsFromCluster(podName, namespace)
	fmt.Println("Fetching metrics for ips:", ips)

	for _, ip := range ips {
		u := "http://" + ip + ":" + podPort + "/metrics"
		r, err := http.Get(u)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("")
			body := r.Body
			bytes, err := ioutil.ReadAll(body)
			if err != nil {
				fmt.Println(err)
			} else {
				auxiliar = string(bytes)

				for _, line := range strings.Split(strings.TrimSuffix(auxiliar, "\n"), "\n") {
					if line[0] != '#' {
						array = strings.Split(line, " ")

						if len(array) == 2 {
							if strings.Contains(line, "{") {
								line = array[0][:len(array[0])-1] + ",nodeip=\"" + ip + "\"} " + array[1]
							} else {
								line = array[0] + "{nodeip=\"" + ip + "\"} " + array[1]
							}
						}
					}

					resp += line + "\n"
				}
			}
		}
	}

	w.Write([]byte(resp))
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	podName = os.Getenv("MTAIL_EXPORTER_POD_NAME")
	namespace = os.Getenv("MTAIL_EXPORTER_NAMESPACE")
	podPort = os.Getenv("MTAIL_EXPORTER_POD_PORT")

	if podName == "" {
		podName = "mtail"
		fmt.Println("No pod name defined, using ", podName)
	}
	if namespace == "" {
		namespace = "default"
		fmt.Println("No pod name defined, using ", namespace)
	}
	if podPort == "" {
		podPort = "3903"
		fmt.Println("No pod port defined, using ", podPort)
	}

	http.HandleFunc("/metrics", exporterHandler)

	port := os.Getenv("MTAIL_EXPORTER_PORT")
	if port == "" {
		port = "8888"
	}
	fmt.Printf("Fetching from pod '%s' port '%s' on namespace '%s'\n", podName, podPort, namespace)

	fmt.Printf("Listening at port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
