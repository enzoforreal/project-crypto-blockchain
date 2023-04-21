package network

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

// Server : La structure du serveur
type Server struct {
	NodeAddress     string          // L'adresse du noeud
	Neighbours      map[string]bool // Les adresses des noeuds voisins
	MessageReceived chan Message    // Le canal de réception de messages
}

// Start : Lance le serveur
func (server *Server) Start() {
	listener, err := net.Listen("tcp", server.NodeAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Écoute des connexions entrantes sur %s\n", server.NodeAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Nouvelle connexion entrante de %s\n", conn.RemoteAddr().String())

		go server.handleConnection(conn)
	}
}

// SendMessage : Envoie un message à un noeud voisin
func (server *Server) SendMessage(message Message, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(message)
	if err != nil {
		return err
	}

	return nil
}

// handleConnection : Gère une nouvelle connexion entrante
func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	var message Message
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&message)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Nouveau message reçu de %s (%s)\n", conn.RemoteAddr().String(), message.Type)

	server.MessageReceived <- message
}

// BroadcastMessage : Diffuse un message à tous les noeuds voisins
func (server *Server) BroadcastMessage(message Message) {
	log.Printf("Diffusion du message (%s) à tous les voisins\n", message.Type)

	encodedMessage, err := gobEncodeMessage(message)
	if err != nil {
		log.Println(err)
		return
	}

	for address := range server.Neighbours {
		err = server.SendMessage(Message{Type: message.Type, Payload: encodedMessage}, address)
		if err != nil {
			log.Printf("Erreur lors de la diffusion du message à %s : %s\n", address, err)
		}
	}
}

// gobEncodeMessage : Encodage Gob pour les messages
func gobEncodeMessage(msg Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(msg)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return buf.Bytes(), nil
}
