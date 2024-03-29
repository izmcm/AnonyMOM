package main

import (
	"cryptoTest"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

var addr = flag.String("addr", "localhost:8084", "http service address")

//------------------------------KEY para lista tal --------------------------
var keyBlack = []byte("0123426789012345")

//--------------------------------------------------------------------------
func makePost() {

	/*formData := url.Values{
		"token":   {"1234567890"},
		"queue":   {"kk"},
		"content": {"vai te tomar no olho do cu biroliro"}, //tentando isso
	}
	*/
	tm1 := time.Now()
	msg := tm1.String() + "se vc nao for so vc nao vai"
	if encrypted, err := cryptoTest.Encrypt(keyBlack, msg); err != nil {
		log.Println(err)
	} else {
		log.Printf("CIPHER KEY: %s\n", keyBlack)
		log.Printf("ENCRYPTED: %s\n", encrypted)

		if decrypted, err := cryptoTest.Decrypt(keyBlack, encrypted); err != nil {
			log.Println(err)
		} else {
			log.Printf("DECRYPTED: %s\n", decrypted)
		}

		// forma o cabeçalho da mensagem ecrypted

		//time inicial

		formData := url.Values{
			"token":   {"1234567890"},
			"queue":   {"filacrypto"},
			"content": {encrypted},
			"action":  {"InsertData"},
		}
		log.Printf("enviu: %s\n", msg)
		//envia a menssagem
		resp, err := http.PostForm("http://localhost:8081", formData)
		if err != nil {
			log.Fatalln(err)
			log.Printf("aqui fudeo")
		}

		log.Println("response made")
		log.Println(resp.Body)
	}
}

func createQueue(token string, queue string, tp int) {
	typeString := strconv.Itoa(tp)
	formData := url.Values{
		"token":  {token},
		"queue":  {queue},
		"type":   {typeString},
		"action": {"CreateQueue"},
	}

	resp, err := http.PostForm("http://localhost:8081", formData)
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
	createQueue("1234567890", "filacrypto", 1)

	log.Printf("making post")
	for i := 0; i < 3; i++ {
		makePost()
	}
}
