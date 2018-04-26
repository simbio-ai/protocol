package protocol

import (
	"bufio"
	"encoding/gob"
	"fmt"
)

//Update is a struct used by nodes in the gossip communication
//It contains all messages that nodes require in order to reach consensus (transactions, votes etc)
type Update struct {
	//Public key of the owning node
	publicKey string
	//Payload containing owning node's messages and round number
	payload *updatePayload
	//Signature obtained by serializing the payload, hashing and encripting with the node's private key
	signature string
}

//Part of the update
type updatePayload struct {
	roundNumber uint
	messages    []Message
}

//String implements Stringer interface for updatePayload structure
func (up *updatePayload) String() string {
	return fmt.Sprintf("%#v", up)
	//or in JSON format json.Marshal(up)
}

//SendUpdate writes the update to the rw *ReadWriter passed as argument
func (u *Update) SendUpdate(rw *bufio.ReadWriter) error {
	//Send update
	enc := gob.NewEncoder(rw)
	err := enc.Encode(u)
	if err != nil {
		return fmt.Errorf("encoding %v: %v", u, err)
	}
	err = rw.Flush()
	if err != nil {
		return fmt.Errorf("flushing %v: %v", rw, err)
	}
	return nil
}

//ReceiveUpdate function tries to read an update from the
//rw *ReadWriter passed as argument and return the result
func ReceiveUpdate(rw *bufio.ReadWriter) (*Update, error) {
	var u Update
	//Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&u)
	if err != nil {
		return nil, fmt.Errorf("decoding handshake request %v: %v", dec, err)
	}
	return &u, nil
}
