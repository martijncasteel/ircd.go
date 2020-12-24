package ircd

import (
	"net"
)

type channel struct {
	name    string
	members map[net.Conn]*client

	topic string
	modes []string
}

// Channel creates a channel and returns it
func Channel(name string) *channel {
	return &channel{
		name:    name,
		members: make(map[net.Conn]*client),

		modes: []string{"t", "n"},
	}
}

func (channel *channel) broadcast(server *server, client *client, message string, args ...interface{}) {
	if inList(channel.modes, "n") && channel.member(client) == nil {
		client.message(ERR_CANNOTSENDTOCHAN, server.name, client.nickname, channel.name)
		return
	}

	for _, recipient := range channel.members {
		if recipient != client {
			recipient.message(message, args...)
		}
	}
}

func (channel *channel) join(server *server, client *client) {

	if inList(channel.modes, "i") {
		client.message(ERR_INVITEONLYCHAN, server.name, client.nickname, channel.name)
		return
	}

	channel.members[client.connection] = client

	if channel.topic != "" {
		client.message(RPL_TOPIC, server.name, client.nickname, channel.name, channel.topic)
	} else {
		client.message(RPL_NOTOPIC, server.name, client.nickname, channel.name)
	}
}

func (channel *channel) part(server *server, client *client) {
	if channel.members[client.connection] == nil {
		client.message(ERR_NOTONCHANNEL, server.name, client.nickname, channel.name)
		return
	}

	// TODO part message? Not part of rfc1459?
	delete(channel.members, client.connection)
}

func (channel *channel) member(client *client) *client {
	for _, member := range channel.members {
		if member == client {
			return client
		}
	}

	return nil
}

func (channel *channel) disconnect(connection net.Conn) {
	delete(channel.members, connection)
}
