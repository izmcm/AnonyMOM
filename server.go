package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
)

func main() {

}

func receiveData(receiveChannel <-chan amqp.Delivery, responseChannel *amqp.Channel) {
}

func postData(data Response, responseChannel *amqp.Channel, clientId int) {

}

func createQueue(name string, channel *amqp.Channel) amqp.Queue {

}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
