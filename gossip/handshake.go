package protocol

import (
	"fmt"
)

//HandShakeMessage is the first message exchanged by nodes
//It contains all data necessary in order to engage in gossip
type HandShakeMessage struct {
	//Node's public key
	publicKey string
	//Payload containing node's current round number and machine state
	payload *handShakePayload
	//Signature obtained by serializing the payload, hashing and encripting with the node's private key
	signature string
}

//Part of HandShakeMessage
type handShakePayload struct {
	roundNumber  uint
	machineState string
}

//String implements Stringer interface for handShakePayload structure
func (hsp *handShakePayload) String() string {
	return fmt.Sprintf("%#v", hsp)
	//or in JSON format json.Marshal(hsp)
}
