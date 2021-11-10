package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-redis/redis"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var redisClient *redis.Client

func main() {
	var ro redis.Options
	ro.Addr = os.Getenv("REDIS_ADDR")

	redisClient = redis.NewClient(&ro)
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis, cannot proceed", err)
	}
	log.Println("👍🏻 Connected to Redis server..")

	kubernetesServiceHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config := new(rest.Config)
	if len(kubernetesServiceHost) > 0 {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}

	// Initialise kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// TODO: potentially watch other resources, i'm only interested in namespaces for now
	watcher, err := clientset.CoreV1().Namespaces().Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for event := range watcher.ResultChan() {
		ns := event.Object.(*v1.Namespace)
		nsJson, err := json.Marshal(ns)
		if err != nil {
			fmt.Println(err)
		}

		switch event.Type {
		case watch.Added:
			fmt.Println("Namespace", ns.ObjectMeta.Name, "added")
			err = publishMessage("namespace added", nsJson)
		case watch.Modified:
			fmt.Printf("Namespace %s modified", ns.ObjectMeta.Name)
			err = publishMessage("namespace modified", nsJson)
		case watch.Deleted:
			fmt.Printf("Namespace %s deleted", ns.ObjectMeta.Name)
			err = publishMessage("namespace deleted", nsJson)
		}

		if err != nil {
			fmt.Println(err)
		}
	}
}

func publishMessage(whatHappened string, data []byte) error {
	fmt.Println(string(data))

	return redisClient.XAdd(&redis.XAddArgs{
		Stream:       "kubernetes",
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"whatHappened": whatHappened,
			"k8sObject":    data,
		},
	}).Err()
}
