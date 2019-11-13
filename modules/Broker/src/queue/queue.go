package main

import (
	"container/list"
	"fmt"
	"message"
)

type AnonyQueue struct {
	name     string
	messages *list.List
}

func (queue *AnonyQueue) pushMessageToQueue(msg message.AnonyMessage) bool {
	queue.messages.PushBack(msg)
	return true
}

func (queue *AnonyQueue) getMessage() message.AnonyMessage {
	element := queue.messages.Front()
	queue.messages.Remove(element)

	return element.Value.(message.AnonyMessage)
}

func createQueue(name string) AnonyQueue {
	queue := AnonyQueue{name: name, messages: list.New()}
	return queue
}

func main() {
	msg := message.AnonyMessage{SenderToken: "men", Queue: "kukuk", Content: "oi amor, to ligando pra saber se foi divertido me trair"}
	msg1 := message.AnonyMessage{SenderToken: "men2", Queue: "2kukuk", Content: "2oi amor, to ligando pra saber se foi divertido me trair"}

	queue := createQueue("kkk")
	queue.pushMessageToQueue(msg)
	queue.pushMessageToQueue(msg1)
	for queue.messages.Len() > 0 {
		e := queue.getMessage()

		fmt.Print(e)
	}
}
