package queue

import (
	"errors"
	"fmt"
	"io/ioutil"
	"message"
	"os"
	"strings"
)

// type AnonyQueue struct {
// 	Name     string
// 	Messages *list.List
// }

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func PushMessageToQueue(msg message.AnonyMessage) bool {
	f, err := os.OpenFile("../database/"+msg.Queue+".txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return false
	}

	newLine := msg.SenderToken + ";" + msg.Content
	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return false
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func GetMessage(name string) (string, error) {
	data, err := ioutil.ReadFile("../database/" + name + ".txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	sliceData := strings.Split(string(data), "\n")
	msg := sliceData[0]

	totalData := strings.Join(sliceData[1:], "\n")
	err = ioutil.WriteFile("../database/"+name+".txt", []byte(totalData), 0644)
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

func CreateQueue(name string) error {
	if !fileExists("../database/" + name + ".txt") {
		_, err := os.Create("../database/" + name + ".txt")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// func main() {
// 	msg := message.AnonyMessage{SenderToken: "men", Queue: "kkk", Content: "oi amor, to ligando pra saber se foi divertido me trair"}

// 	createQueue("kkk")
// 	createQueue("kk1")

// 	pushMessageToQueue(msg)
// 	msgGet := getMessage("kkk")
// 	fmt.Println(msgGet)
// }
