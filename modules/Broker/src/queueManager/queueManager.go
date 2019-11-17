package main

import (
	"message"
	"queue"
)

type AnonyQueueManager struct {
	pool     map[string]queue.AnonyQueue
	userPool map[string]int
}

func (manager *AnonyQueueManager) insertMessageToQueue(content message.AnonyMessage) bool {
	err := queue.CreateQueue(content.Queue)
	if err != nil {
		return false
	} else {
		ok := queue.PushMessageToQueue(content)
		return ok
	}
}

func (manager *AnonyQueueManager) getMessageFromQueue(queueName string) string {
	msg := queue.GetMessage(queueName)
	return msg
}

func (manager *AnonyQueueManager) subscribeUserToQueue(user string, queueName string) bool {
	// TODO: implement this method so that the user
	// must be authenticated, it is, the user exists
	// so the user can be identified as a person not
	// a robot (?)

	manager.userPool[user] = 1
	return true
}

func (manager *AnonyQueueManager) unsubscribeUserToQueue(user string, queueName string) bool {
	// TODO: implement this method so that the user
	// must be authenticated, it is, the user exists
	// so the user can be identified as a person not
	// a robot (?)

	manager.userPool[user] = 0
	return true
}

func main() {
	manager := AnonyQueueManager{pool: map[string]int{"Mark": 10, "Sandy": 20}, userPool: make(map[string]queue.AnonyQueue)}
	message1 := message.AnonyMessage{SenderToken: "men", Queue: "kkk", Content: "oi amor, to ligando pra saber se foi divertido me trair"}
	// manager.insertMessageToQueue(message)

	message2 := message.AnonyMessage{SenderToken: "123", Queue: "kkk", Content: "oloko meu2"}
	// manager.insertMessageToQueue(message2)

	message3 := message.AnonyMessage{SenderToken: "123", Queue: "kkk", Content: "oloko meu3"}
	// manager.insertMessageToQueue(message3)

	// returnMsg := getMessageFromQueue("kkk")
	// fmt.Fprintln(returnMsg)
}
