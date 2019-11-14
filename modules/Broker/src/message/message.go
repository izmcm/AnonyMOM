package message

type AnonyMessage struct {
	SenderToken string
	Queue       string
	Content     string
}

type CryptedMessage struct {
	PublicKey string
	Message   string
}
