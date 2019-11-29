// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// https://medium.com/@masnun/making-http-requests-in-golang-dd123379efe7

package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	// 	"os"
	// 	"os/signal"
	// 	"time"
	// 	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8082", "http service address")

func makePost() {
	//
	formData := url.Values{
		"token":   {"1234567890"},
		"queue":   {"kk"},
		"content": {"vai te tomar no olho do cu biroliro"},
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
	makePost()
}
