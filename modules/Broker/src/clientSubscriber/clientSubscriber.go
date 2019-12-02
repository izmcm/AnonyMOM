package clientSubscriber

import (
	"container/list"
	"log"
	"net/http"
	"net/url"
	// 	"os"
	// 	"os/signal"
	// 	"time"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	Token         string
	Addrs         string
	Connection    *websocket.Conn
	OutputChannel chan string
}

// must be without the HTTP
func New(addrs string, token string, channel chan string) Subscriber {
	s := Subscriber{Addrs: addrs, Token: token, OutputChannel: channel}
	return s
}

func (s *Subscriber) Subscribe(queues list.List) {
	// Setup a channel to check if there's a command in the keyboard
	// 	interrupt := make(chan os.Signal, 1)
	// 	signal.Notify(interrupt, os.Interrupt)

	// Setup the websocket URL
	u := url.URL{Scheme: "ws", Host: s.Addrs, Path: "/"}
	log.Printf("connecting to %s", u.String())

	// Insert the headers
	header := http.Header{}
	if s.Token != "" {
		header.Add("token", s.Token)
	}
	for el := queues.Front(); el != nil; el = el.Next() {
		header.Add("Subscriptions", el.Value.(string))
	}

	// Dial to the websocket service
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}

	s.Connection = c

	// This is a channel where the message will be set when we receive it from the server
	done := make(chan struct{})
	go func(s *Subscriber, done chan struct{}) {
		defer close(done)
		for {
			// 			log.Println("trying to reach:", s.Connection)
			// Read message from websocket, blocking
			_, message, err := s.Connection.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			// log.Printf("recv: %s", message)
			// Blocking action to put the message data somewhere else
			if s.OutputChannel != nil {
				msg := string(message)
				s.OutputChannel <- msg
			} else {
				log.Println(s.OutputChannel)
			}
		}
	}(s, done)
}
