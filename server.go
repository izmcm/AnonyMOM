package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"sync"
)

var amqpURL string = "amqp://guest:guest@localhost:5672/"
var wg sync.WaitGroup

var amqpQueues map[string]amqp.Queue = make(map[string]amqp.Queue)

type AddTask struct {
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

	// create mainQueue
	mainQueue := createQueue("MAIN", amqpChannel)

	// the server will deliver that many messages to consumers before acknowledgments are received
	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	// create consumer to receive messages from mainQueue
	mainChannel, err := amqpChannel.Consume(
		mainQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	// receive data from queue
	wg.Add(1)
	go receiveData(mainChannel, amqpChannel)
	wg.Wait()
}

// Function to read data from a channel and call another function to send response
// Params:
//  - receiveChannel: channel from which messages will be read
//  - responseChannel: channel to which messages will be send
func receiveData(receiveChannel <-chan amqp.Delivery, responseChannel *amqp.Channel) {
	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())

		for message := range receiveChannel {
			log.Printf("Received a message: %s", message.Body)

			addTask := &AddTask{}

			err := json.Unmarshal(message.Body, addTask)
			handleError(err, "Error decoding JSON")

			num := addTask.Number
			cid := addTask.ClientId
			data := Response{Number: num}

			// return the data to channel
			go postData(data, responseChannel, cid)

			if err := message.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	// stop for program termination
	<-stopChan
	wg.Done()
}

// Function to post data in a channel
// Params:
//  - data: data that will be post
//  - responseChannel: channel to which messages will be send
//  - clientId: client id that initially sent message
func postData(data Response, responseChannel *amqp.Channel, clientId int) {
	body, err := json.Marshal(data)
	handleError(err, "Error encoding JSON")

	// the ansQueue will be named with your clientId
	// so, the answer only be disponible to the server and the right client
	channelName := fmt.Sprintf("Ans%d", clientId)
	if _, ok := amqpQueues[channelName]; ok {

	} else {
		amqpQueues[channelName] = createQueue(channelName, responseChannel)
	}

	// publish message in queue
	err = responseChannel.Publish("", channelName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	handleError(err, "Error publishing message")
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
