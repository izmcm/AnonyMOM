// TODO: maybe create an engine to allow the host to make a post when you are
// into the websocket

// TODO: Implement whitelist
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
	"net/url"
	"queue"
	"queueManager"
	"strconv"
	"strings"
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
var numOfHostsInQueueChannel map[string]chan int = make(map[string]chan int) // store a list o channels
var numOfHostsInQueue map[string]int = make(map[string]int)                  // store a list o channels
var numOfHosts chan int = make(chan int)
var addr = flag.String("addr", "localhost:8082", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var manager queueManager.AnonyQueueManager = queueManager.AnonyQueueManager{Pool: make(map[string]*queue.AnonyQueue), UserPool: make(map[string]int)}
var broker Broker = Broker{Broker_id: "0", Manager: manager}
var runningQueues list.List

/**
 * Treats the data from the HTTP request handler
 **/
func (broker Broker) GETandPOST(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	fmt.Println("Received a request")

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
		fmt.Println(r.Header["token"])

		// Verify if we are asking for the websocket upgrade
		upgradeToWebsocket := false
		for _, h := range r.Header["Upgrade"] {
			if h == "websocket" {
				upgradeToWebsocket = true
				break
			}
		}

		// if we should upgrade then we upgrade and generate a websocket
		if upgradeToWebsocket {
			var queuesNotAllowed list.List

			// TODO: send to the user it's automatic generated token
			connection := SetupWebsocketConnection(w, r)
			host := Host{ID: RandStringBytes(15), Connection: connection}
			if len(r.Header["token"]) > 0 {
				host.ID = r.Header["token"][0]
			}
			fmt.Println("--------------")

			for _, queueName := range r.Header["Subscriptions"] {

				// Check if the queue exists before let the user subscribe
				q, err := broker.Manager.GetQueueNamed(queueName)
				if err != nil {
					queuesNotAllowed.PushBack(queueName)
					fmt.Println("queue named", queueName, "doesn't exist!")
					continue
				}

				// Verify if the user can register in the queue before running it
				canSub := broker.Manager.CheckUserRights(host.ID, q, 0)
				if !canSub {
					queuesNotAllowed.PushBack(queueName)
					fmt.Println("queue named", queueName, "doesn't exist!")
					continue
				}

				fmt.Println("subscribed to:", queueName)
				var shouldRun bool = true

				// Check if there's a specific actor for look at this queue
				for queueElement := runningQueues.Front(); queueElement != nil && shouldRun; queueElement = queueElement.Next() {
					if queueElement.Value.(string) == queueName {
						fmt.Println("found queue with the same name here")
						fmt.Println("1: ", queueElement.Value.(string))
						fmt.Println("2: ", queueName)
						shouldRun = false
					}
				}

				// If we need start the broker broadcast for this specific queue
				if shouldRun {
					// Setup a channel for the queue and start a listenner
					fmt.Println("Registering listenner and channel for queue", queueName)
					numOfHostsInQueueChannel[queueName] = make(chan int)
					go broker.listenAndBroadcastQueue(queueName)
				}

				host.Subscriptions.PushBack(queueName)
				fmt.Println("Prev in queue <", queueName, ">:", numOfHostsInQueue[queueName])
				numOfHostsInQueue[queueName] += 1
				fmt.Println("Hosts in queue <", queueName, ">:", numOfHostsInQueue[queueName])
			}
			hostList.PushBack(host)
			for _, q := range r.Header["Subscriptions"] {
				if c, ok := numOfHostsInQueueChannel[q]; ok {
					c <- numOfHostsInQueue[q]
					fmt.Println()
				}
			}
		} else {
			w.Write([]byte("Request Not processed by the server\n"))
		}
	case "POST":
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		data := r.Form
		fmt.Println("POST HERE")
		fmt.Println(data)
		fmt.Println("")

		// Check if it's a valid action
		action := SelectAction(data)

		//-------------------------------------------------
		// Create a new queue
		//-------------------------------------------------
		if action == "CreateQueue" && data["type"] != nil {
			tp, err := strconv.Atoi(data["type"][0])
			if err != nil {
				fmt.Println("Type", tp, "not allowed")
				w.Write([]byte("Queue Cannot be registered, type not allowed"))
			} else {
				// queue := queue.AnonyQueue{Name: data["queue"][0], Owner: data["token"][0], Type: tp}
				queue := queue.New(data["queue"][0], data["token"][0], tp, 0)
				_, err := broker.Manager.GetQueueNamed(queue.Name)

				if err != nil {
					broker.Manager.RegisterQueue(&queue)
					fmt.Println("Queue", queue.Name, "registered for user", queue.Owner)
					w.Write([]byte("Queue Registered"))
				} else {
					fmt.Println("Queue", queue.Name, "cannot be registered for user", queue.Owner)
					w.Write([]byte("Queue Cannot be registered"))
				}
			}
		}

		//-------------------------------------------------
		// Insert a message into a queue
		//-------------------------------------------------
		if action == "InsertData" && data["content"] != nil {
			message := message.AnonyMessage{SenderToken: data["token"][0], Queue: data["queue"][0], Content: data["content"][0]}
			status, err := broker.Manager.InsertMessageToQueue(message)

			if err != nil {
				fmt.Println("Error inserting data to queue\n")
				w.Write([]byte("Error inserting data to queue\n"))
			} else if status {
				fmt.Println("Data inserted\n")
				w.Write([]byte("Data inserted\n"))
			} else {
				fmt.Println("You are not allowed to insert data in this queue")
				w.Write([]byte("You are not allowed to insert data in this queue"))
			}
		}

		if action == "BlackList" && data["List"] != nil {
			q, err := broker.Manager.GetQueueNamed(data["queue"][0])
			if err != nil {
				fmt.Println("Error blacklisting users\n")
				w.Write([]byte("Error blacklisting users\n"))
				return
			}

			if q.Owner != data["token"][0] {
				fmt.Println("You don't own this queue\n")
				fmt.Println("the owner is", q.Owner, "not you,", data["token"][0])
				w.Write([]byte("You don't own this queue\n"))
				return
			}

			if q.Type != 1 {
				fmt.Println("Queue doesn't support blacklist\n")
				w.Write([]byte("Queue doesn't support blacklist\n"))
				return
			}

			fmt.Println("--------------------")
			for _, user := range data["List"] {
				fmt.Println("Blacklisting user", user)
				q.BlackList.PushBack(user)
			}
			fmt.Println("starting serialization")
			// TODO: new goroutine here
			q.SerializeQueue()
			fmt.Println("serialization finished")
			fmt.Println("--------------------")

		}

		// Request not processed
		if action == "" {
			w.Write([]byte("Request Not processed by the server\n"))
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

// Select the action that the user wants to do
func SelectAction(data url.Values) string {
	action := ""

	if data["action"] != nil && data["token"] != nil && data["queue"] != nil {
		action = data["action"][0]
	}

	return action
}

// make a simple SetupWebsocketConnection server for the websocket service
func SetupWebsocketConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return nil
	}
	c.WriteMessage(websocket.TextMessage, []byte("bora dale boy"))
	c.WriteMessage(websocket.TextMessage, []byte("tais registrado boy!"))

	return c
}

func (broker Broker) listenAndBroadcastQueue(queue string) {
	hostNum := 0
	c := numOfHostsInQueueChannel[queue]
	for {
		select {
		case data := <-c:
			hostNum = data
		default:
			if hostNum != 0 {
				message, err := broker.Manager.GetMessageFromQueue(queue)
				if err != nil {
					// No messages
					if err.Error() != "Empty queue" {
						fmt.Println(err)
					}
				} else {
					fmt.Println("Making broadcast")
					fmt.Println()
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

			newMessageBytes := []byte(message)
			// set maximum TCP network packet size (1460 bytes)
			bytesSize := len([]byte(message))
			if bytesSize < 1460 {
				trash := strings.Repeat("0", 1459-bytesSize)
				newMessage := message + ";" + trash
				newMessageBytes = []byte(newMessage)
			}

			err := c.WriteMessage(websocket.TextMessage, newMessageBytes)

			if err != nil {
				log.Println("write:", err)
				// WARNING: removing from list, check if cause problems in loop
				numOfHostsInQueue[queue] -= 1
				numOfHostsInQueueChannel[queue] <- numOfHostsInQueue[queue]
				hostList.Remove(host)
				continue
			}
		}
	}
}

// Setup a timer to broadcast a time based message
func broadcastTimer() {
	fmt.Println("Setting up the timer broadcast")
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

	fmt.Println("Setting up the HTTP Handler")

	// manager := queueManager.AnonyQueueManager{UserPool: make(map[string]int)}
	// broker := Broker{Broker_id: "0", Manager: manager}
	http.HandleFunc("/", broker.GETandPOST)

	// go broadcastTimer()
	// go broker.listenAndBroadcastQueue("kk", numOfHosts)

	// http.ListenAndServe(":3001", nil)
	http.ListenAndServe(":8082", nil)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
