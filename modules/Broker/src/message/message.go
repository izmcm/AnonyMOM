package message

type AnonyMessage struct {
	SenderToken string
	Queue       string
	Content     string
}

// not use this
type CryptedMessage struct {
	PublicKey string
	Message   string
}
