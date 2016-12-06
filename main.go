package main

import (
	//Import necessary external packages
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/watch"

	"k8s.io/client-go/rest"
)

func main() {
	// creates the in-cluster config

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	eventInterface, err := clientset.Core().Events("").Watch(v1.ListOptions{})
	eventCh := eventInterface.ResultChan()
	if err != nil {
		panic(err.Error())
	}

	for {
		select {
		case event := <-eventCh:
			ref, err := api.GetReference(event.Object)
			if err == nil {
				switch event.Type {
				case watch.Added:
					fmt.Printf("Added: %v\n", ref.Kind);
				case watch.Modified:
					fmt.Printf("Modified: %v\n", ref.Kind);
				case watch.Deleted:
					fmt.Printf("Deleted: %v\n", ref.Kind);
				case watch.Error:
					fmt.Printf("Error: %v\n", ref.Kind);
				}
			}else{
				fmt.Printf("%v %v\n", event.Type, err);
			}

		}
	}

}
