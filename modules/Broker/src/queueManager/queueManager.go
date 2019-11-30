package queueManager

import (
	"container/list"
	// "errors"
	"message"
	"queue"
	// "strconv"
)

// TODO: Setup the type queue so we can define metadata about it such as if the
// queue has a whitelist or a blacklist and then define the functions into the
// queue for the access token for the blacklist or whitelist or even if it's
// a public queue and anyone can write into it.

type AnonyQueueManager struct {
	Pool     map[string]queue.AnonyQueue
	UserPool map[string]int
}

// Global constants
const READ_OPERATION int = 1
const WRITE_OPERATION int = 1

func (manager *AnonyQueueManager) RegisterQueue(meta queue.AnonyQueue) (bool, error) {
	_, err := meta.SerializeQueue()
	// _, err := meta.CreateQueue()
	manager.Pool[meta.Name] = meta

	if err != nil {
		return false, err
	}

	return true, nil
}

func (manager *AnonyQueueManager) InsertMessageToQueue(msg message.AnonyMessage) (bool, error) {
	q, err := manager.GetQueueNamed(msg.Queue)
	if err != nil {
		return false, err
	}

	// check if the user can write in this queue
	canWrite := manager.CheckUserRights(msg.SenderToken, q, WRITE_OPERATION)

	if canWrite {
		err, ok := queue.PushMessageToQueue(msg)
		return ok, err
	}

	return false, nil
}

// func (manager *AnonyQueueManager) GetMessageFromQueue(queueName string, userToken string) (string, error) {
// 	q, err := manager.GetQueueNamed(queueName)
// 	if err != nil {
// 		return "", err
// 	}

// 	// check if the user can write in this queue
// 	canRead := manager.CheckUserRights(userToken, q, READ_OPERATION)

// 	if canRead {
// 		msg, err := queue.GetMessage(queueName)
// 		return msg, err
// 	}

// 	// TODO: check if the user has the possibility to read from this queue
// 	return "", errors.New("User not authorized to read from the queue")
// }

func (manager *AnonyQueueManager) GetMessageFromQueue(queueName string) (string, error) {
	_, err := manager.GetQueueNamed(queueName)
	if err != nil {
		return "", err
	}

	msg, err := queue.GetMessage(queueName)
	return msg, err
}

func (manager *AnonyQueueManager) SubscribeUserToQueue(user string, queueName string) bool {
	// TODO: Check if the user can subscribe to this queue,
	// how we are working with the ability to the user read
	// only when it's subscribed maybe, but just maybe the
	// we don't need to implement it
	manager.UserPool[user] = 1
	return true
}

func (manager *AnonyQueueManager) UnsubscribeUserToQueue(user string, queueName string) bool {
	// TODO: Check if the user can subscribe to this queue,
	// how we are working with the ability to the user read
	// only when it's subscribed maybe, but just maybe the
	// we don't need to implement it
	manager.UserPool[user] = 0
	return true
}

func (manager *AnonyQueueManager) GetQueueNamed(queueName string) (queue.AnonyQueue, error) {
	if queue, ok := manager.Pool[queueName]; ok {
		return queue, nil
	}

	queue, err := queue.ReadQueueFromFile(queueName)
	if err != nil {
		return queue, err
	}

	return queue, nil
}

func ContainsString(lst list.List, value string) bool {
	isIn := false
	for el := lst.Front(); el != nil && !isIn; el.Next() {
		val := el.Value.(string)
		if value == val {
			isIn = false
		}
	}
	return isIn
}

// Operation:
// 0 - read
// 1 - write
func (manager *AnonyQueueManager) CheckUserRights(userToken string, q queue.AnonyQueue, operation int) bool {
	allow := true
	if operation == 0 {
		if q.Type == 1 {
			allow = !ContainsString(q.BlackList, userToken)
		} else if q.Type == 2 {
			allow = ContainsString(q.WhiteList, userToken)
		}
	} else if operation == 1 {
		if q.Type == 1 {
			allow = !ContainsString(q.WritersBlackList, userToken)
		} else if q.Type == 2 {
			allow = ContainsString(q.WritersWhiteList, userToken)
		}
	}

	return allow
}

// func main() {
// 	// manager := AnonyQueueManager{pool: map[string]int{"Mark": 10, "Sandy": 20}, userPool: make(map[string]queue.AnonyQueue)}
// 	manager := AnonyQueueManager{userPool: make(map[string]int)}
// 	message1 := message.AnonyMessage{SenderToken: "men", Queue: "kkk", Content: "oi amor, to ligando pra saber se foi divertido me trair"}
// 	// manager.insertMessageToQueue(message)

// 	message2 := message.AnonyMessage{SenderToken: "123", Queue: "kkk", Content: "oloko meu2"}
// 	// manager.insertMessageToQueue(message2)

// 	message3 := message.AnonyMessage{SenderToken: "123", Queue: "kkk", Content: "oloko meu3"}
// 	// manager.insertMessageToQueue(message3)

// 	// returnMsg := getMessageFromQueue("kkk")
// 	// fmt.Fprintln(returnMsg)
// }
