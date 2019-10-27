package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"os"
	"sync"
)

var amqpURL string = "amqp://guest:guest@localhost:5672/"
var wg sync.WaitGroup

type MessageToSend struct {
	Number   int
	ClientId int
}

type Response struct {
	Number int
}

func main() {
	// connect to RabbitMQ
	conn, err := amqp.Dial(amqpURL)
	handleError(err, "Can't connect to AMQP")
	defer conn.Close()

	// create RabbitMQ channel
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")
	defer amqpChannel.Close()

	// create clientId
	clientId := rand.Intn(99999)
	channelName := fmt.Sprintf("Ans%d", clientId)

	// create queues
	mainQueue := createQueue("MAIN", amqpChannel)

	// the ansQueue will be named with your clientId
	// so, the answer only be disponible to the server and the right client
	ansQueue := createQueue(channelName, amqpChannel)

	// the server will deliver that many messages to consumers before acknowledgments are received
	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	// create consumer to receive messages from ansQueue
	ansChannel, err := amqpChannel.Consume(
		ansQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	// post in mainQueue and receive data in ansChannel (create with ansQueue)
	postData(amqpChannel, mainQueue, clientId)
	wg.Add(1)
	go receiveData(ansChannel)
	wg.Wait()
}

// Function to post data in a channel
// Params:
//  - data: data that will be post
//  - responseChannel: channel to which messages will be send
//  - clientId: client id that initially sent message
func postData(channel *amqp.Channel, queue amqp.Queue, clientId int) {
	message := MessageToSend{Number: rand.Intn(999), ClientId: clientId}

	body, err := json.Marshal(message)
	handleError(err, "Error encoding JSON")

	// publish message in queue
	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	handleError(err, "Error publishing message")
}

// Function to read data from a channel
// Params:
//  - receiveChannel: channel from which messages will be read
func receiveData(receiveChannel <-chan amqp.Delivery) {
	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())

		for message := range receiveChannel {
			log.Printf("Received a message: %s", message.Body)

			res := &Response{}
			err := json.Unmarshal(message.Body, res)

			handleError(err, "Error decoding JSON")

			if err := message.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	// Stop for program termination
	<-stopChan
	wg.Done()
}

// Function to create a queue from a channel
// Params:
//  - name: queue's name
//  - channel: channel that will be created the queue
// Return:
//  - the queue that was created
func createQueue(name string, channel *amqp.Channel) amqp.Queue {
	queue, err := channel.QueueDeclare(name, true, false, false, false, nil)
	handleError(err, fmt.Sprintf("Could not declare %s queue", name))
	return queue
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
