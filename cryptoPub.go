package main

import (
	"clientPublisher"
	"cryptoTest"
	"fmt"
)

var keyBlack = []byte("0123426789012345")

func main() {
	pub := clientPublisher.New("http://localhost:8082", "1234567890")
	fmt.Println(pub)
	pub.CreateQueue("novafila", 1)
	if encrypted, err := cryptoTest.Encrypt(keyBlack, "pokebola vai"); err == nil {
		pub.Publish("novafila", encrypted)
	}
	if encrypted, err := cryptoTest.Encrypt(keyBlack, "pokemon 123"); err == nil {
		pub.Publish("novafila", encrypted)
	}
	if encrypted, err := cryptoTest.Encrypt(keyBlack, "é nois que voa bruxão"); err == nil {
		pub.Publish("novafila", encrypted)
	}
	// pub.Publish("benchmarkQueues", anonymize.AnonymizeMessage("pokebola vai"))
	// pub.Publish("novafila", anonymize.AnonymizeMessage("pokemon 123"))
	// pub.Publish("novafila", anonymize.AnonymizeMessage("é nois que voa bruxão"))
}
