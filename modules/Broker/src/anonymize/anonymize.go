package anonymize

import (
	// 	"fmt"
	"strings"
)

func AnonymizeMessage(message string) string {
	newMessage := message

	// set maximum TCP network packet size (1460 bytes)
	bytesSize := len([]byte(message))
	if bytesSize < 1460 {
		trash := strings.Repeat("0", 1459-bytesSize)
		newMessage = message + ";" + trash
	}

	return newMessage
}

func DeanonymizeMessage(message string) string {
	parts := strings.Split(message, ";")
	parts = parts[:len(parts)-1]
	deanonMsg := strings.Join(parts[:], ";")
	return deanonMsg
}

// func main() {
// 	a := anonymizeMessage("oi porra")
// 	fmt.Println(a)
// 	b := deanonymizeMessage(a)
// 	fmt.Println(b)
// }
