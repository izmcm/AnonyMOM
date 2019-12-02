package main

import (
	"anonymize"
	"clientPublisher"
	"fmt"
)

func main() {
	pub := clientPublisher.New("http://localhost:8082", "1234567890")
	fmt.Println(pub)
	pub.CreateQueue("novafila", 1)
	pub.Publish("novafila", anonymize.AnonymizeMessage("pokebola vai"))
	pub.Publish("novafila", anonymize.AnonymizeMessage("pokemon 123"))
	pub.Publish("novafila", anonymize.AnonymizeMessage("é nois que voa bruxão"))
}
