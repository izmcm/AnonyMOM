package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
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

func Blacklist(token string, queue string, user string) {
	formData := url.Values{
		"token":  {token},
		"queue":  {queue},
		"List":   {user},
		"action": {"BlackList"},
	}

	resp, err := http.PostForm("http://localhost:8082", formData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("response made")
	log.Println(resp.Body)
}

func createQueue(token string, queue string, tp int) {
	typeString := strconv.Itoa(tp)
	formData := url.Values{
		"token":  {token},
		"queue":  {queue},
		"type":   {typeString},
		"action": {"CreateQueue"},
	}

	resp, err := http.PostForm("http://localhost:8082", formData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("response made")
	log.Println(resp.Body)
}

func insertData(token string, queue string, content string) {
	formData := url.Values{
		"token":   {token},
		"queue":   {queue},
		"content": {content},
		"action":  {"InsertData"},
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

	log.Printf("Creating queue")
	createQueue("1234567890", "newqueue", 1)

	Blacklist("1234567890", "novafila", "pokemon")
	insertData("1234567890", "novafila", RandStringBytes(20))
	insertData("pokemon", "novafila", RandStringBytes(20))
	insertData("1234567890", "novafila", RandStringBytes(20))

	// log.Printf("making post")
	// for i := 0; i < 30; i++ {
	// 	insertData("1234567890", "novafila", RandStringBytes(20))
	// }
}
