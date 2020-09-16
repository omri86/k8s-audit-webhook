package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/apiserver/pkg/apis/audit/v1"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		var events v1.EventList
		err = json.Unmarshal(body, &events)
		if err != nil {
			http.Error(w, "failed to unmarshal audit events", http.StatusBadRequest)
			return
		}
		// Iterate and filter audit events
		for _, event := range events.Items {
			if isPodCreation(event) {
				fmt.Printf("Pod creation event detected: %+v\n", event)
			}
		}
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// isPodCreation returns true if the given event is of a pod creation
func isPodCreation(event v1.Event) bool {
	return event.Verb == "create" &&
		event.Stage == v1.StageResponseComplete &&
		event.ObjectRef != nil &&
		event.ObjectRef.Resource == "pods"
}
