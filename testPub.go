package main

import (
	"clientPublisher"
	"fmt"
)

func main() {
	pub := clientPublisher.New("http://localhost:8082", "1234567890")
	fmt.Println(pub)
	pub.CreateQueue("novafila", 1)
	pub.Publish("novafila", "pokebola vai")
	pub.Publish("novafila", "pokemon 123")
	pub.Publish("novafila", "é nois que voa bruxão")
}
