package clientPublisher

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Publisher struct {
	Token string
	Addrs string
}

// addrs must be http address
func New(addrs string, token string) Publisher {
	p := Publisher{Token: token, Addrs: addrs}
	return p
}

func (p *Publisher) Publish(queue string, content string) {
	formData := url.Values{
		"token":   {p.Token},
		"queue":   {queue},
		"content": {content},
		"action":  {"InsertData"},
	}

	resp, err := http.PostForm(p.Addrs, formData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("response made")
	log.Println(resp.Body)
}

func (p *Publisher) CreateQueue(queue string, tp int) {
	typeString := strconv.Itoa(tp)
	formData := url.Values{
		"token":  {p.Token},
		"queue":  {queue},
		"type":   {typeString},
		"action": {"CreateQueue"},
	}

	resp, err := http.PostForm(p.Addrs, formData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("response made")
	log.Println(resp.Body)
}

func (p *Publisher) Blacklist(queue string, userToBlackList string) {
	formData := url.Values{
		"token":  {p.Token},
		"queue":  {queue},
		"List":   {userToBlackList},
		"action": {"BlackList"},
	}

	resp, err := http.PostForm(p.Addrs, formData)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("response made")
	log.Println(resp.Body)
	// TODO: add a return
}
