package ircd

import (
	"log"
	"net"
)

type server struct {
	name   string
	config *config

	channels map[string]*channel
	clients  map[net.Conn]*client
}

// Server is called to create a chat server
func Server(config *config) *server {
	server := server{
		name:     "martijncasteel.com",
		config:   config,
		channels: make(map[string]*channel),
		clients:  make(map[net.Conn]*client),
	}

	server.channels["#general"] = Channel("#general")
	return &server
}

func (server *server) Run() {
	listener, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatalf("Unable to start the server: %s", err.Error())
		return
	}

	defer listener.Close()
	log.Printf("Started server on :8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Unable to accept connection: %s", err.Error())
			continue
		}

		go Client(server, conn)
	}
}

func (server *server) disconnect(connection net.Conn) {
	client := server.clients[connection]

	// if client has called quit, this will be called again
	if client == nil {
		return
	}

	log.Printf("%s (%s) disconnected", client.nickname, connection.RemoteAddr().String())

	for _, channel := range server.channels {
		channel.disconnect(connection)
	}

	connection.Close()
	delete(server.clients, connection)
}

func (server *server) channel(name string) *channel {
	for key, value := range server.channels {
		if key == name {
			return value
		}
	}

	return nil
}

func (server *server) client(name string) *client {
	for _, client := range server.clients {
		if client.nickname == name {
			return client
		}
	}

	return nil
}

func (server *server) login(client *client, username string, password string) bool {
	//TODO password and config etc
	if !inList(client.modes, "o") {
		client.modes = append(client.modes, "o")
		return true
	}

	return false
}

func inList(list []string, str string) bool {
	for _, val := range list {
		if val == str {
			return true
		}
	}

	return false
}
