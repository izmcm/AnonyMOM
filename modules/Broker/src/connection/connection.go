package main

import (
	"container/list"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"message"
	"net/http"
	"queueManager"
	// "strconv"
	"time"
)

type Broker struct {
	Broker_id string
	Manager   queueManager.AnonyQueueManager
}

type Host struct {
	ID            string
	Connection    *websocket.Conn
	Subscriptions list.List
}

var hostList list.List

func (broker Broker) GETandPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		// TODO: Setup the data to run into another goroutine
		// in order to process the request, but for now let's
		// just do the stuff like if we are in a single goroutine
		// and fuck this shit!

		fmt.Println(r)
		fmt.Println(r.URL)
		fmt.Println(r.Header)

		fmt.Println(r.Header["Subscriptions"])

		// Verify if we are asking for the websocket upgrade
		if r.Header["Upgrade"][0] == "websocket" {
			connection := echo(w, r)
			// TODO: Extract to a host connection function later
			// TODO: Setup the userID in the Header
			host := Host{ID: "abc", Connection: connection}
			hostList.PushBack(host)
		} else {
			w.Write([]byte("Request Not processed by the server\n"))
		}
	case "POST":
		// TODO: verify who can put data in which channel before just
		// allowing the users to insert any data, but for now, fuck you!
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		data := r.Form
		fmt.Println("POST HERE")
		fmt.Println(data)
		fmt.Println("")

		if data["token"] != nil && data["queue"] != nil && data["content"] != nil {
			message := message.AnonyMessage{SenderToken: data["token"][0], Queue: data["queue"][0], Content: data["content"][0]}
			broker.Manager.InsertMessageToQueue(message)
			w.Write([]byte("Data inserted\n"))
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
	// fmt.Println(c)
	c.WriteMessage(websocket.TextMessage, []byte("bora dale boy"))
	c.WriteMessage(websocket.TextMessage, []byte("tais registrado boy!"))

	return c
}

func (broker Broker) listenAndBroadcastQueue(name string) {
	for {
		message, err := broker.Manager.GetMessageFromQueue(name)
		if err != nil {
			// No messages
		} else {
			// TODO: make the broadcast only to the queue

			// TODO: check if there's someone connected before
			// remove from the queue
			fmt.Println("Making broadcast")
			broadcastMessage(message)
		}
	}
}

func broadcastMessage(message string) {
	id := 0
	for host := hostList.Front(); host != nil; host = host.Next() {
		id += 1
		c := host.Value.(Host).Connection
		err := c.WriteMessage(websocket.TextMessage, []byte(message))
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
		case <-ticker.C:
			// fmt.Println(t)
			broadcastMessage("tick")
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
	// go broker.listenAndBroadcastQueue("kk")

	// http.ListenAndServe(":3001", nil)
	http.ListenAndServe(":8082", nil)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
