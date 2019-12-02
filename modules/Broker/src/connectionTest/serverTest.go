package main

import (
	"connection"
	"fmt"
)

func main() {
	fmt.Println("Setting up the HTTP Handler")
	broker := connection.New("1", "localhost:8082")
	broker.ListenConnections()
}
