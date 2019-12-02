package main

import (
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
		// 		select {
		// 		case msg := <-channel:
		// 			fmt.Println("received the message: ", msg)
		// 		default:
		// 			//
		// 		}
		// fmt.Println("waiting for message from channel")
		msg := <-channel
		fmt.Println("received message: ", msg)
	}
}
