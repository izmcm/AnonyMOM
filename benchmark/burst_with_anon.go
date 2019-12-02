package main

import (
	"clientPublisher"
	"clientSubscriber"
	"container/list"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var messageLimit int = 14000
var numOfPublishers int = 2

func main() {
	numOfSubs := 1
	pub := clientPublisher.New("http://localhost:8082", "benchMarker0")
	pub.CreateQueue("benchmarkQueue", 1)

	var publishers list.List
	publishers.PushBack(pub)

	wg.Add(numOfPublishers + numOfSubs)

	// create publisher
	for i := 1; i < numOfPublishers+1; i++ {
		p := clientPublisher.New("http://localhost:8082", "benchMarker"+strconv.Itoa(i))
		publishers.PushBack(p)
		publishMessages(p)
	}

	// create a listenner
	channel := make(chan string)
	sub := clientSubscriber.New("localhost:8082", "kkk", channel)

	var queues list.List
	queues.PushBack("benchmarkQueue")

	sub.Subscribe(queues)
	// 	go listennerManager("benchmark0", channel)
	go listennerManager("benchmark_with_"+strconv.Itoa(numOfPublishers)+"_publishers.csv", channel)

	// wait until the threads finish
	wg.Wait()
}

func publishMessages(pub clientPublisher.Publisher) {
	// publish everything to here
	for i := 0; i < messageLimit; i++ {
		tm := time.Now().UnixNano()
		str := strconv.FormatInt(tm, 10)
		pub.Publish("benchmarkQueue", str)
	}
	wg.Done()
}

func listennerManager(filename string, channel chan string) {
	ct := 0
	f, err := os.OpenFile("./data/"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	for {
		// Check to stop the program when needed
		if ct >= messageLimit+2 {
			wg.Done()
			break
		}
		ct += 1

		msg := <-channel

		parts := strings.Split(msg, ";")
		parts = parts[len(parts)-1:]

		tm, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			continue
		}
		dur := time.Now().UnixNano() - tm
		if _, err := f.WriteString(strconv.FormatInt(dur, 10) + "\n"); err != nil {
			//
		}
	}
}
