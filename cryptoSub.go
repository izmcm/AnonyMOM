package main

import (
	"clientSubscriber"
	"container/list"
	"cryptoTest"
	"fmt"
	"log"
	"strings"
)

var keyBlack = []byte("0123426789012345")

func main() {
	channel := make(chan string)

	sub := clientSubscriber.New("localhost:8082", "dale boy", channel)
	fmt.Println(sub)

	var queues list.List
	queues.PushBack("novafila")

	// It start some goroutines so must lock the software to work properly

	sub.Subscribe(queues)

	for {
		msg := <-channel
		fmt.Println("received message: ", msg)

		parts := strings.Split(msg, ";")
		parts = parts[len(parts)-1:]
		data := parts[0]

		if decrypted, err := cryptoTest.Decrypt(keyBlack, data); err != nil {
			log.Println(err)
		} else {
			log.Println(decrypted)
		}

	}
}
