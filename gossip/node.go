package protocol

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

const (
	//Port is a constant specifying the default listening port of the server
	Port = ":61000"
)

//Node is a struct emulating a network node
type Node struct {
	publicKey   string
	privateKey  *rsa.PrivateKey
	isListening bool
	//port string
	listener net.Listener
	network  []*Node
}

//NewNode creates and returns a new Node struct object
func NewNode() *Node {
	return &Node{"needs implementing", nil, false, nil, nil}
}

//GetCurrentMachineState returns the current machine state i.e. the latest block hash
func (node *Node) GetCurrentMachineState() string {
	return "needs implementing"
}

//GetCurrentMessages returns a slice of messages collected in the current round,
//such as transactions received from wallet apps, votes etc.
func (node *Node) GetCurrentMessages() []Message {
	return nil //needs implementing
}

//GetCurrentRoundNumber returns the current round number
func (node *Node) GetCurrentRoundNumber() uint {
	return 0 //needs implementing
}

//Hash method computes the hash of 'str' argument and retuns a string containing
//the hex representation
func (node *Node) Hash(str string) string {
	sum := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", sum)
}

//SignMessage returns a signature of the node corresponding to the strucutre passed as argument
func (node *Node) SignMessage(msg fmt.Stringer) (string, error) {
	rng := rand.Reader
	message := []byte(msg.String())
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rng, node.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", signature), nil
}

//MakeHandshake returns a HandShakeMessage to be sent to another node
func (node *Node) MakeHandshake() (*HandShakeMessage, error) {
	hsp := &handShakePayload{node.GetCurrentRoundNumber(), node.GetCurrentMachineState()}
	sig, err := node.SignMessage(hsp)
	return &HandShakeMessage{node.publicKey, hsp, sig}, err
}

//MakeCurrentUpdate creates and returns the update for the current round
func (node *Node) MakeCurrentUpdate() (*Update, error) {
	upl := &updatePayload{node.GetCurrentRoundNumber(), node.GetCurrentMessages()}
	sig, err := node.SignMessage(upl)
	return &Update{node.publicKey, upl, sig}, err
}

//SendHandshake creates a handshake message for the current node and writes it to the
//rw *ReadWriter passed as argument
func (node *Node) SendHandshake(rw *bufio.ReadWriter) error {
	//Create handshakeMessage
	hsResponse, err := node.MakeHandshake()
	if err != nil {
		return fmt.Errorf("creating handshake %v: %v", node, err)
	}
	//Send handshake
	enc := gob.NewEncoder(rw)
	err = enc.Encode(hsResponse)
	if err != nil {
		return fmt.Errorf("encoding %v: %v", hsResponse, err)
	}
	err = rw.Flush()
	if err != nil {
		return fmt.Errorf("flushing %v: %v", rw, err)
	}
	return nil
}

//ReceiveHandshake method tries to read a handshake message from the
//rw *ReadWriter passed as argument and return the result
func (node *Node) ReceiveHandshake(rw *bufio.ReadWriter) (*HandShakeMessage, error) {
	var hsReceived HandShakeMessage
	//Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&hsReceived)
	if err != nil {
		return nil, fmt.Errorf("decoding handshake request %v: %v", dec, err)
	}
	return &hsReceived, nil
}

//Listen runs a loop and accepts connections from other processes
func (node *Node) Listen() error {
	var err error
	//Try to start listening
	node.listener, err = net.Listen("tcp", Port)
	//Handle potential errors
	if err != nil {
		return fmt.Errorf("Unable to listen on %v: %v", node.listener.Addr(), err)
	}
	//Start loop to handle incoming connections
	node.isListening = true
	for node.isListening {
		//Try accepting a connection request.
		conn, err := node.listener.Accept()
		//Handle potential errors
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		//Handle incoming connections
		go node.handleIncoming(conn)
	}
	return nil
}

//handleIncoming talks to other nodes and receives/sends updates
func (node *Node) handleIncoming(conn net.Conn) error {
	//Close connection when surrounding function returns
	//As an optimisation, it might be better to use the same connection for multiple/random number of rounds
	defer conn.Close()
	//Create a bufio.ReadWriter from the connection
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	//Receive handshake
	hsRequest, err := node.ReceiveHandshake(rw)
	if err != nil {
		return fmt.Errorf("receiving handshake request %v: %v", rw, err)
	}
	//Create response handshake and send it to the client
	err = node.SendHandshake(rw)
	if err != nil {
		return fmt.Errorf("sending handshake response %v: %v", rw, err)
	}
	//TODO:
	//Compare machine states and reconcile differencess if necessary
	if node.GetCurrentMachineState() != hsRequest.payload.machineState {
		//TODO:Reconcile
		//recurse on states until the first occurance of a different state
		//compare nodes counted in consensus for said round
		//for any missing/exclded nodes, each part (client and server) must present 50%+1 votes of excluding in the round
		//the server will exclude from consensus all necessary nodes
	}

	//Read update requests and respond by sending updates
	//Optimised behavior must include requests/response for parts of updates
	for {
		//TODO: Receive/read update requests
		//TODO: compare current round updates
		//TODO: if have update send update //myUpdate comes first
		//TODO: else send request for update
		//TODO: receive  update
		//TODO: end loop
	}
}

//Gossip runs a loop and talks to other processes
func (node *Node) Gossip() {
	for {
		return
	}
}
