package marshaller

import (
	"encoding/json"
	"miop"
)

type Marshaller struct{}

func (Marshaller) Marshall(msg miop.Packet) []byte {
	r, err := json.Marshal(msg)
	if err != nil {
		// TODO
	}

	return r
}

func (Marshaller) Unmarshall(msg []byte) miop.Packet {
	r := miop.Packet{}
	err := json.Unmarshal(msg, &r)
	if err != nil {
		// TODO
	}

	return r
}
