package main

import (
	"bufio"
	"clientPublisher"
	"clientSubscriber"
	"container/list"
	"fmt"
	"log"
	"os"
	"strings"
	// "time"
)

func main() {
	pub := clientPublisher.New("http://localhost:8082", "1234567890")
	fmt.Println(pub)

	fmt.Print("Escolha um user: ")
	username := getFromKeyboard()
	fmt.Print("Olá " + username + "Escolha uma opção:\n")

	fmt.Print("1. Create room\n2. Enter room\n")
	op := getFromKeyboard()

	fmt.Print("Qual o nome da sala? ")
	room := getFromKeyboard()

	fmt.Println(op)
	if op[0] == []byte("1")[0] {
		createRoom(room, pub)
	}
	_, subscribeChannel := enterRoom(room)

	ch := make(chan string)
	go func(ch chan string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err := reader.ReadString('\n')
			msg = strings.Split(msg, "\n")[0]
			if err != nil {
				close(ch)
				return
			}
			ch <- msg
		}
	}(ch)

stdinloop:
	for {
		select {
		case stdin, ok := <-ch:
			if !ok {
				break stdinloop
			} else {
				postData(stdin, room, pub)
			}
		// case <-time.After(1 * time.Second):
		// 	receiveData(sub)
		case msg := <-subscribeChannel:
			fmt.Println("recebido:", msg)
		}
	}
	fmt.Println("Done, stdin must be closed")
}

func postData(msg string, room string, pub clientPublisher.Publisher) {
	fmt.Println("post: " + msg)
	pub.Publish(room, msg)
}

func receiveData(sub clientSubscriber.Subscriber, channel chan string) {

}

func enterRoom(room string) (clientSubscriber.Subscriber, chan string) {
	channel := make(chan string)
	sub := clientSubscriber.New("localhost:8082", room, channel)
	var queues list.List
	queues.PushBack(room)
	sub.Subscribe(queues)

	return sub, channel
}

func createRoom(name string, pub clientPublisher.Publisher) {
	pub.CreateQueue(name, 1)
}

func getFromKeyboard() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.Split(text, "\n")[0]
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
