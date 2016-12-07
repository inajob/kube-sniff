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

func printEvent(e *v1.Event) {
	/*
	fmt.Printf("Name: %v\n", e.Name)
	fmt.Printf("Message: %v\n", e.Message)
	fmt.Printf("Reason: %v\n", e.Reason)
	fmt.Printf("Source: %v\n", e.Source)
	fmt.Printf("InvolvedObject: %v\n", e.InvolvedObject)
	fmt.Printf("  Kind: %v\n", e.InvolvedObject.Kind)
	fmt.Printf("  Name: %v\n", e.InvolvedObject.Name)
        */
	fmt.Printf("Reason: %v\n", e.Reason)
	fmt.Printf("  Kind: %v\n", e.InvolvedObject.Kind)
	fmt.Printf("  Name: %v\n", e.InvolvedObject.Name)
	fmt.Println();
}

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
	eventInterface, err := clientset.CoreV1().Events("").Watch(v1.ListOptions{})
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
					fmt.Println("## Added:")
					printEvent(ev)
				case watch.Modified:
					fmt.Println("## Modified:")
					printEvent(ev)
				case watch.Deleted:
					fmt.Println("## Deleted:")
					printEvent(ev)
				case watch.Error:
					fmt.Println("## Error:")
					printEvent(ev)
				}

			} else{
				fmt.Printf("error\n")
			}

		}
	}

}
