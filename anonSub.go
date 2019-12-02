package main

import (
	"anonymize"
	"clientSubscriber"
	"container/list"
	"fmt"
)

func main() {
	channel := make(chan string)

	sub := clientSubscriber.New("localhost:8082", "kkk", channel)
	fmt.Println(sub)

	var queues list.List
	queues.PushBack("novafila")

	// It start some goroutines so must lock the software to work properly

	sub.Subscribe(queues)

	for {
		msg := <-channel
		fmt.Println("received message: ", anonymize.DeanonymizeMessage(msg))
	}
}
