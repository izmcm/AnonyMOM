package main

import (
	"fmt"
	// "io/ioutil"
	"container/list"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"message"
	"net/http"
	"queueManager"
	"strconv"
	"time"
)

type Broker struct {
	Broker_id string
	Manager   queueManager.AnonyQueueManager
}

type Host struct {
	ID         string
	Connection *websocket.Conn
}

var hostList list.List

// var connections = make(map[websocket.Conn]*list.List) // host:privilégios

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

		// TODO: we still must make the checkup for the handshake
		// of a websocket and also make a proper usage of the
		// websocket communication method.
		//
		// Verify if we are asking for the websocket upgrade
		if r.Header["Upgrade"][0] == "websocket" {
			connection := echo(w, r)
			// TODO: Extract to a host connection function later
			host := Host{ID: "abc", Connection: connection}
			hostList.PushBack(host)
		} else if params["token"] != nil && params["queue"] != nil && params["content"] != nil {
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

var addr = flag.String("addr", "localhost:8082", "http service address")
var upgrader = websocket.Upgrader{} // use default options

// TODO: make the server send data for the client after stabilish
// the connection in the websocket, also must define which of the
// queues each client can have access at all.
// and every single update into a queue must be logged into the
// listenners of the websocket with a notification of new data
// so the client requests back the data from the queue (?) or
// just chooses a subscriber from a pool or a broadcast value
// for the subscribers to put the data.

// make a simple echo server for the websocket service
func echo(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return nil
	}
	// defer c.Close()
	// for {
	// 	mt, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// 	fmt.Println(mt)
	// 	err = c.WriteMessage(mt, message)
	// 	// err = c.WriteMessage(mt, []byte("kk bora"))
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		break
	// 	}
	// }

	// connections[c].PushBack()

	c.WriteMessage(websocket.TextMessage, []byte("bora dale boy"))
	c.WriteMessage(websocket.TextMessage, []byte("tais registrado boy!"))

	return c
}

func broadcastMessage() {
	id := 0
	for host := hostList.Front(); host != nil; host = host.Next() {
		id += 1
		c := host.Value.(Host).Connection
		err := c.WriteMessage(websocket.TextMessage, []byte("hello "+strconv.Itoa(id)))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func broadcastTimer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			fmt.Println(t)
			broadcastMessage()
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	manager := queueManager.AnonyQueueManager{UserPool: make(map[string]int)}
	broker := Broker{Broker_id: "0", Manager: manager}
	http.HandleFunc("/", broker.GETandPOST)
	// http.HandleFunc("/echo", echo)

	go broadcastTimer()

	// http.ListenAndServe(":3001", nil)
	http.ListenAndServe(":8082", nil)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
