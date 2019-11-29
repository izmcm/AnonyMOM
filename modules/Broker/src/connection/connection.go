// TODO: lookup for active queues in the broker and
// setup the broadcast to every single queue

// TODO: send data from the queue if there's a new user

// TODO: maybe create an engine to make a post when you are into the websocket
package main

import (
	"container/list"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"message"
	"net/http"
	"queueManager"
	// "strconv"
	"time"
)

// useful functions
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Data types
type Broker struct {
	Broker_id string
	Manager   queueManager.AnonyQueueManager
}

type Host struct {
	ID            string
	Connection    *websocket.Conn
	Subscriptions list.List
}

// Global variables
var hostList list.List
var numOfHostsInQueue list.List // store a list o channels
var numOfHosts chan int = make(chan int)
var addr = flag.String("addr", "localhost:8082", "http service address")
var upgrader = websocket.Upgrader{} // use default options

/**
 * Treats the data from the HTTP request handler
 **/
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
		// TODO: verify if the upgrade to websocket exists
		// in a bigger list of fields
		if r.Header["Upgrade"][0] == "websocket" {
			connection := echo(w, r)
			// TODO: Extract to a host connection function later
			// TODO: Setup the userID in the Header
			host := Host{ID: RandStringBytes(15), Connection: connection}
			fmt.Println("--------------")
			// TODO: define if every host can subscribe in every
			// queue
			for _, q := range r.Header["Subscriptions"] {
				fmt.Println("subscribed to:", q)
				host.Subscriptions.PushBack(q)
			}
			hostList.PushBack(host)
			// TODO: understand how to make it non blocking or not, don't know
			// yet
			numOfHosts <- hostList.Len()
			// Update the number of hosts in the channel
			// select {
			// case numOfHosts <- hostList.Len():
			// 	fmt.Println("queue size updated")
			// default:
			// 	fmt.Println("error updating queue")
			// }
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

func (broker Broker) listenAndBroadcastQueue(queue string, c chan int) {
	hostNum := 0
	for {
		select {
		case data := <-c:
			hostNum = data
		default:
			// TODO: Modify to check the num of hosts for queue other
			// than the total number of hosts.
			if hostNum != 0 {
				message, err := broker.Manager.GetMessageFromQueue(queue)
				if err != nil {
					// No messages
				} else {
					// TODO: make the broadcast only to the queue
					fmt.Println("Making broadcast")
					broadcastMessage(message, queue)
				}
			}
		}
	}
}

// Send a message to all the listenners
func broadcastMessage(message string, queue string) {
	id := 0
	isMember := false
	for host := hostList.Front(); host != nil; host = host.Next() {
		isMember = false
		// TODO: update this id value
		id += 1
		subs := host.Value.(Host).Subscriptions
		// Check if host is subscribed in this queue
		for el := subs.Front(); el != nil; el = el.Next() {
			q := el.Value.(string)
			if q == queue {
				isMember = true
				break
			}
		}
		// If registered to the queue then send the message to them
		if isMember {
			c := host.Value.(Host).Connection
			err := c.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				// WARNING: removing from list, check if cause problems in loop
				hostList.Remove(host)
				continue
			}
		}
	}
}

// Setup a timer to broadcast a time based message
func broadcastTimer() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// fmt.Println(t)
			broadcastMessage("tick", "kk")
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	manager := queueManager.AnonyQueueManager{UserPool: make(map[string]int)}
	broker := Broker{Broker_id: "0", Manager: manager}
	http.HandleFunc("/", broker.GETandPOST)

	go broadcastTimer()
	go broker.listenAndBroadcastQueue("kk", numOfHosts)

	// http.ListenAndServe(":3001", nil)
	http.ListenAndServe(":8082", nil)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
