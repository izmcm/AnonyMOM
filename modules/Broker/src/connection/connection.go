package main

import (
	"fmt"
	// "io/ioutil"
	// "log"
	"message"
	"net/http"
	"queueManager"
)

type Broker struct {
	Broker_id string
	Manager   queueManager.AnonyQueueManager
}

func (broker Broker) GETandPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		for k, v := range r.URL.Query() {
			fmt.Printf("%s: %s\n", k, v)
		}
		// w.Write([]byte("Received a GET request\n"))
		// TODO: Setup the data to run into another goroutine
		// in order to process the request, but for now let's
		// just do the stuff like if we are in a single goroutine
		// and fuck this shit!
		fmt.Println(r)
		fmt.Println(r.URL)
		params := r.URL.Query()
		fmt.Println(params["user"])

		if params["token"] != nil && params["queue"] != nil && params["content"] != nil {
			// TODO: decrypt and make sense of request
			message := message.AnonyMessage{SenderToken: params["token"][0], Queue: params["queue"][0], Content: params["content"][0]}
			broker.Manager.InsertMessageToQueue(message)
			w.Write([]byte("Data inserted\n"))
		} else if params["token"] != nil || params["queue"] != nil {
			// TODO: decrypt and make sense of request
			// message := message.AnonyMessage{SenderToken: params["token"][0], Queue: params["queue"][0], Content: params["content"][0]}
			message := params["queue"][0]
			messageData := broker.Manager.GetMessageFromQueue(message)
			w.Write([]byte(messageData + "\n"))
		} else {
			w.Write([]byte("Request Not processed by the server\n"))
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}

}

func main() {
	manager := queueManager.AnonyQueueManager{UserPool: make(map[string]int)}
	broker := Broker{Broker_id: "0", Manager: manager}
	http.HandleFunc("/", broker.GETandPOST)
	http.ListenAndServe(":3001", nil)
}
