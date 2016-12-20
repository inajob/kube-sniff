package main

import (
	//Import necessary external packages
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
	"golang.org/x/net/websocket"
	"net/http"
)

func printEvent(e *v1.Event) {
	fmt.Printf("Reason: %v\n", e.Reason)
	fmt.Printf("  Kind: %v\n", e.InvolvedObject.Kind)
	fmt.Printf("  Name: %v\n", e.InvolvedObject.Name)
	fmt.Println();
}

func emitEvent(e *v1.Event, msgCh chan v1.Event) {
	msgCh <- *e
}


func msgHandler(ws *websocket.Conn, msgCh chan v1.Event){
	for {
		s := <- msgCh
		if err := websocket.JSON.Send(ws, s); err != nil {
			fmt.Println("cant send ", s);
			return
		}
	}
}

func startServer(mc chan v1.Event) {
	http.Handle("/msg", websocket.Handler(func(ws *websocket.Conn){msgHandler(ws, mc)}))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	if err := http.ListenAndServe(":9999", nil); err != nil {
		panic("ListenAndServe: " + err.Error())
	}
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

	msgCh := make(chan v1.Event);
	
	eventInterface, err := clientset.CoreV1().Events("").Watch(v1.ListOptions{})
	eventCh := eventInterface.ResultChan()
	if err != nil {
		panic(err.Error())
	}

	go startServer(msgCh)

	for {
		select {
		case event := <-eventCh:
			ev, ok := event.Object.(*v1.Event)
			if ok {
				switch event.Type {
				case watch.Added:
					fmt.Println("## Added:")
					printEvent(ev)
					emitEvent(ev, msgCh)
				case watch.Modified:
					fmt.Println("## Modified:")
					printEvent(ev)
					emitEvent(ev, msgCh)
				case watch.Deleted:
					fmt.Println("## Deleted:")
					printEvent(ev)
					emitEvent(ev, msgCh)
				case watch.Error:
					fmt.Println("## Error:")
					printEvent(ev)
					emitEvent(ev, msgCh)
				}
			} else{
				//fmt.Printf("error\n")
			}

		}
	}

}
