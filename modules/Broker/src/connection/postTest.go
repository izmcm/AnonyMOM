package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	// 	"os"
	// 	"os/signal"
	// 	"time"
	// 	"github.com/gorilla/websocket"
)

// useful functions
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var addr = flag.String("addr", "localhost:8082", "http service address")

func makePost() {
	// 	formData := url.Values{
	// 		"token":   {"1234567890"},
	// 		"queue":   {"kk"},
	// 		"content": {"vai te tomar no olho do cu biroliro"},
	// 	}
	formData := url.Values{
		"token":   {"1234567890"},
		"queue":   {"novafila"},
		"content": {RandStringBytes(20)},
	}

	resp, err := http.PostForm("http://localhost:8082", formData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("response made")
	log.Println(resp.Body)
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	// 	interrupt := make(chan os.Signal, 1)
	// 	signal.Notify(interrupt, os.Interrupt)

	// u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("making post")
	for i := 0; i < 20; i++ {
		makePost()
	}
}
