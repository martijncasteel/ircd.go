package ircd

import (
	"log"
	"net"
)

type server struct {
	name   string
	config *Config

	channels map[string]*channel
	clients  map[net.Conn]*client
}

// Server is called to create a chat server
func Server(config *Config) *server {
	server := server{
		name:     config.Name,
		config:   config,
		channels: make(map[string]*channel),
		clients:  make(map[net.Conn]*client),
	}

	for _, element := range config.Channels {
		// server.channels["#general"] = Channel("#general")
		server.channels[element] = Channel(element)
	}
	return &server
}

func (server *server) Run() {
	listener, err := net.Listen("tcp", server.config.Address)

	if err != nil {
		log.Fatalf("Unable to start the server: %s", err.Error())
		return
	}

	defer listener.Close()
	log.Printf("Started server on %s", server.config.Address)

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

	log.Printf("- %s disconnected (%s)", connection.RemoteAddr().String(), client.nickname)

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

// Error checks for error and logs it
func Error(err error) bool {
	if err != nil {
		log.Print(err)
		return true
	}
	return false
}
