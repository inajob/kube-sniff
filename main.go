package main

import (
	//Import necessary external packages
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	//"k8s.io/client-go/pkg/api/meta"
	//"k8s.io/client-go/pkg/api"
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
			/*
			obj, err := meta.Accessor(event.Object);
			if err == nil {
				switch event.Type {
				case watch.Added:
					fmt.Printf("Added: %v\n", obj.GetName());
				case watch.Modified:
					fmt.Printf("Modified: %v\n", obj.GetName());
				case watch.Deleted:
					fmt.Printf("Deleted: %v\n", obj.GetName());
				case watch.Error:
					fmt.Printf("Error: %v\n", obj.GetName());
				}
			}else{
				fmt.Printf("%v %v\n", event.Type, err);
			}*/
			ev, ok := event.Object.(*v1.Event)
			if ok {
				//fmt.Printf("Name: %v, Message: %v", ev.Name, ev.Message);
				switch event.Type {
				case watch.Added:
					fmt.Printf("Added: %v %v\n", ev.Name, ev.Message);
				case watch.Modified:
					fmt.Printf("Modified: %v %v\n", ev.Name, ev.Message);
				case watch.Deleted:
					fmt.Printf("Deleted: %v %v\n", ev.Name, ev.Message);
				case watch.Error:
					fmt.Printf("Error: %v %v\n", ev.Name, ev.Message);
				}

			} else{
				fmt.Printf("error\n");
			}

		}
	}

}
