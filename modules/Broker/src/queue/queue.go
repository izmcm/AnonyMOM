package queue

import (
	"container/list"
	"errors"
	"fmt"
	"io/ioutil"
	// "message"
	"os"
	// "queueManager"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Metadata about a specific queue
// TODO: implement the persistence of this
// metadata into some files or databases.
// Types:
// 1 - public: anyone can read and write, but the ones in the blacklist cannot
// read or write, and the ones in the writtersBlackList cannot write
// 2 - private: only the ones in the whitelist can read and only the ones in the
// Writters whitelist can write
type AnonyQueue struct {
	Name             string
	Owner            string
	Type             int
	WhiteList        list.List
	BlackList        list.List
	WritersWhiteList list.List
	WritersBlackList list.List
	Mux              *sync.Mutex
	Throughput       int64
	LastHit          time.Time
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func New(name string, owner string, tp int, throughput int64) AnonyQueue {
	return AnonyQueue{Name: name, Owner: owner, Type: tp, Throughput: throughput, LastHit: time.Now()}
}

func deleteFile(path string) error {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// TODO: implement the system to make the queue persist it's metadata
// func CreateQueue(name string, owner string, queueType string) (string, error) {
func (queue *AnonyQueue) CreateQueue() (string, error) {
	var name string = queue.Name
	var owner string = queue.Owner
	var queueType string = strconv.Itoa(queue.Type)

	if queueType != "2" && queueType != "1" {
		return "", errors.New("Unknown queue Type")
	}
	if !fileExists("./database/" + name + ".txt") {
		_, err := os.Create("./database/" + name + ".txt")
		if err != nil {
			fmt.Println(err)
			return "Erro na criação da fila", err
		}

		_, err = os.Create("./meta/" + name + ".txt")
		if err != nil {
			fmt.Println(err)
			return "Erro na criação dos metadados da fila", err
		}

		fmt.Println(queueType)
		metaData := owner + "\n" + queueType
		err = ioutil.WriteFile("./meta/"+name+".txt", []byte(metaData), 0644)
		if err != nil {
			fmt.Println("File error", err)
			return "Erro na criação dos metadados da fila", err
		}
	} else {
		err := errors.New("This queue already exists")
		return "Fila já existente", err
	}

	return "Sucesso", nil
}

func (q *AnonyQueue) SerializeQueue() (string, error) {

	if q.Type < 1 || q.Type > 2 {
		return "Type" + strconv.Itoa(q.Type) + " not allowed", errors.New("Unknown queue Type")
	}

	// if !fileExists("./database/" + q.Name + ".txt") {
	if !fileExists("./meta/"+q.Name+".txt") || !fileExists("./database/"+q.Name+".txt") {
		_, err := os.Create("./database/" + q.Name + ".txt")
		if err != nil {
			fmt.Println(err)
			return "Erro na criação da fila", err
		}

		_, err = os.Create("./meta/" + q.Name + ".txt")
		if err != nil {
			fmt.Println(err)
			return "Erro na criação dos metadados da fila", err
		}

	} else {
		// err := errors.New("This queue already exists")
		// return "Fila já existente", err
	}

	metaData := q.Owner + "\n" + strconv.Itoa(q.Type)
	if q.Type == 1 {
		metaData = metaData + "\n" + strconv.Itoa(q.BlackList.Len())
		metaData = metaData + "\n" + strconv.Itoa(q.WritersBlackList.Len())

		for el := q.BlackList.Front(); el != nil; el = el.Next() {
			metaData = metaData + "\n" + el.Value.(string)
		}
		for el := q.WritersBlackList.Front(); el != nil; el = el.Next() {
			metaData = metaData + "\n" + el.Value.(string)
		}
	} else if q.Type == 2 {
		metaData = metaData + "\n" + strconv.Itoa(q.WhiteList.Len())
		metaData = metaData + "\n" + strconv.Itoa(q.WritersWhiteList.Len())

		for el := q.WhiteList.Front(); el != nil; el = el.Next() {
			metaData = metaData + "\n" + el.Value.(string)
		}

		for el := q.WritersWhiteList.Front(); el != nil; el = el.Next() {
			metaData = metaData + "\n" + el.Value.(string)
		}
	}

	err := ioutil.WriteFile("./meta/"+q.Name+".txt", []byte(metaData), 0644)
	if err != nil {
		fmt.Println("File error", err)
		return "Erro na criação dos metadados da fila", err
	}

	return "Sucesso", nil
}

func ReadQueueFromFile(queueName string) (AnonyQueue, error) {
	// Split the data by lines
	data, err := ioutil.ReadFile("./meta/" + queueName + ".txt")
	if err != nil {
		// queue := AnonyQueue{Name: "", Owner: "", Type: -1}
		queue := New("", "", -1, 0)
		return queue, err
	}
	sliceData := strings.Split(string(data), "\n")
	name := queueName
	owner := sliceData[0]
	tp, _ := strconv.Atoi(sliceData[1])

	// queue := AnonyQueue{Name: name, Owner: owner, Type: tp}
	queue := New(name, owner, tp, 0)

	if tp == 1 {
		blackListSize, _ := strconv.Atoi(sliceData[2])
		writersBlackListSize, _ := strconv.Atoi(sliceData[3])
		for i := 4; i < blackListSize; i++ {
			queue.BlackList.PushBack(sliceData[i])
		}
		for i := 4 + blackListSize; i < writersBlackListSize; i++ {
			queue.WritersBlackList.PushBack(sliceData[i])
		}
	} else if tp == 2 {
		whiteListSize, _ := strconv.Atoi(sliceData[2])
		writersWhiteListSize, _ := strconv.Atoi(sliceData[3])
		for i := 4; i < whiteListSize; i++ {
			queue.WhiteList.PushBack(sliceData[i])
		}
		for i := 4 + whiteListSize; i < writersWhiteListSize; i++ {
			queue.WritersWhiteList.PushBack(sliceData[i])
		}
	}

	return queue, nil
}
