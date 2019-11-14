package main

import (
	"container/list"
	"fmt"
	"message"
	"queue"
)

// type AnonyQueue struct {
// 	name     string
// 	messages *list.List
// }
type AnonyQueueManager struct {
	AnonyQueue *list.List
}

func (manager *AnonyQueueManager) getQueueNamed(name string) { //AnonyQueue {
	fmt.Print(maneger.AnonyQueue.len())
}

func (manager *AnonyQueueManager) insertMessageToQueue(message AnonyMessage) bool {
	queue := manager.getQueueNamed(message.name)
	queue.pushMessageToQueue(message)
	return true
}

func (manager *AnonyQueueManager) subscribeUserToQueue(user string, queue string) {

}

func (manager *AnonyQueueManager) unsubscribeUserToQueue(user string, queue string) {

}

func main() {

}
