package queueManager

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	"message"
	"os"
	"queue"
	"strings"
	"sync"
	"time"
)

// TODO: Setup the type queue so we can define metadata about it such as if the
// queue has a whitelist or a blacklist and then define the functions into the
// queue for the access token for the blacklist or whitelist or even if it's
// a public queue and anyone can write into it.

type AnonyQueueManager struct {
	Pool     map[string]*queue.AnonyQueue
	UserPool map[string]int
}

// Global constants
const READ_OPERATION int = 1
const WRITE_OPERATION int = 1

func (manager *AnonyQueueManager) RegisterQueue(meta *queue.AnonyQueue) (bool, error) {
	_, err := meta.SerializeQueue()
	// _, err := meta.CreateQueue()
	if meta.Mux == nil {
		meta.Mux = &sync.Mutex{}
	}
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
	canGo := time.Now().Sub(q.LastHit).Nanoseconds()/1000 > q.Throughput

	if canWrite && canGo {
		err, ok := manager.PushMessageToQueue(msg)
		q.LastHit = time.Now()
		return ok, err
	}

	return false, nil
}

func (manager *AnonyQueueManager) PushMessageToQueue(msg message.AnonyMessage) (error, bool) {
	// lock queue mutex
	q, err := manager.GetQueueNamed(msg.Queue)
	if err != nil {
		return err, false
	}

	// addr := &q.Mux
	// fmt.Printf("push: %p\n", addr)
	q.Mux.Lock()

	f, err := os.OpenFile("../database/"+msg.Queue+".txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err, false
	}

	newLine := msg.SenderToken + ";" + msg.Content
	_, err = fmt.Fprintln(f, newLine)

	// unlock queue mutex
	q.Mux.Unlock()

	if err != nil {
		fmt.Println(err)
		f.Close()
		return err, false
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err, false
	}

	return nil, true
}

func (manager *AnonyQueueManager) GetMessageFromQueue(queueName string) (string, error) {
	// _, err := manager.GetQueueNamed(queueName)
	// if err != nil {
	// 	return "", err
	// }

	msg, err := manager.GetMessage(queueName)
	return msg, err
}

func (manager *AnonyQueueManager) GetMessage(name string) (string, error) {
	// lock queue mutex
	q, err := manager.GetQueueNamed(name)
	if err != nil {
		return "", err
	}

	// addr := &q.Mux
	// fmt.Printf("get: %p\n", addr)
	q.Mux.Lock()

	data, err := ioutil.ReadFile("../database/" + name + ".txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	sliceData := strings.Split(string(data), "\n")
	msg := sliceData[0]

	totalData := strings.Join(sliceData[1:], "\n")
	err = ioutil.WriteFile("../database/"+name+".txt", []byte(totalData), 0644)

	// unlock queue mutex
	q.Mux.Unlock()

	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}

	if msg == "" {
		err := errors.New("Empty queue")
		return msg, err
	}

	return msg, nil
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

func (manager *AnonyQueueManager) GetQueueNamed(queueName string) (*queue.AnonyQueue, error) {
	// already should have a mutex
	if queue, ok := manager.Pool[queueName]; ok {
		if queue.Mux == nil {
			queue.Mux = &sync.Mutex{}
		}
		return queue, nil
	}

	// must create a new mutex
	q, err := queue.ReadQueueFromFile(queueName)
	queue := &q
	if err != nil {
		return queue, err
	}

	// WARNING: colateral effect
	queue.Mux = &sync.Mutex{}
	manager.Pool[queueName] = queue
	return queue, nil
}

func ContainsString(lst *list.List, value string) bool {
	isIn := false
	for el := lst.Front(); el != nil; el = el.Next() {
		if value == el.Value.(string) {
			isIn = true
			break
		}
	}
	return isIn
}

// Operation:
// 0 - read
// 1 - write
func (manager *AnonyQueueManager) CheckUserRights(userToken string, q *queue.AnonyQueue, operation int) bool {
	allow := true
	if operation == 0 {
		if q.Type == 1 {
			allow = !ContainsString(&q.BlackList, userToken)
		} else if q.Type == 2 {
			allow = ContainsString(&q.WhiteList, userToken)
		}
	} else if operation == 1 {
		if q.Type == 1 {
			allow = !ContainsString(&q.BlackList, userToken) && !ContainsString(&q.WritersBlackList, userToken)
		} else if q.Type == 2 {
			allow = ContainsString(&q.WhiteList, userToken) && ContainsString(&q.WritersWhiteList, userToken)
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
