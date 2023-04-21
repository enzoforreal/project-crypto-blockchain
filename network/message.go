package network

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Message : La structure d'un message
type Message struct {
	Type    string // Le type de message (ex: "block", "transaction", etc.)
	Payload []byte // Le contenu du message (sérialisé sous forme de bytes)
}

// GobEncode : Encodage Gob pour les messages
func (msg *Message) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(msg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode : Décodage Gob pour les messages
func (msg *Message) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
