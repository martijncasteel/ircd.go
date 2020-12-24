package ircd

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	connection net.Conn
	registered bool

	realname string
	username string
	nickname string

	modes []string
}

// Client creates a new connected client
func Client(server *server, connection net.Conn) {
	client := client{
		connection: connection,
		modes:      []string{},
	}

	server.clients[connection] = &client
	client.run(server)
}

// prefix <servername> | <nick> [ '!' <user> ] [ '@' <host> ]
func (client *client) prefix(server *server) string {
	if client.username != "" {
		return client.nickname + "!" + client.username
	}

	return client.nickname
}

func (client *client) run(server *server) {

	for {
		packet, err := bufio.NewReader(client.connection).ReadString('\n')
		if err != nil {
			// io.EOF is received for a closed connection
			server.disconnect(client.connection)
			return
		}

		packet = strings.Trim(packet, "\r\n")
		args := strings.SplitN(packet, " ", 2)
		args[0] = strings.ToUpper(args[0])

		if !client.registered {
			if !inList([]string{"PASS", "NICK", "USER", "QUIT"}, args[0]) {
				client.message(ERR_NOTREGISTERED, server.name)
				continue
			}
		}

		client.process(server, args[0], args...)
	}
}

func (client *client) process(server *server, cmd string, args ...string) {

	switch cmd {

	case "NICK":
		client.nick(server, args...)
	case "USER":
		client.user(server, args...)

	case "PRIVMSG":
		client.privmsg(server, args...)
	case "NOTICE":
		client.privmsg(server, args...)

	case "JOIN":
		client.join(server, args...)
	case "PART":
		client.part(server, args...)

	case "MODE":
	case "INVITE":

	case "OPER":
		client.oper(server, args...)
	case "QUIT":
		server.disconnect(client.connection)

	default:
		client.message(ERR_UNKNOWNCOMMAND, server.name, client.nickname, cmd)
	}
}

// :martijn PRIVMSG amber :Hi there
func (client *client) message(message string, args ...interface{}) {
	message = fmt.Sprintf(message, args...)
	client.connection.Write([]byte(message + "\r\n"))
}
