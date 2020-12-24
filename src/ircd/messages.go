package ircd

import (
	"regexp"
	"strings"
)

func (sender *client) nick(server *server, args ...string) {

	if len(args) < 2 {
		sender.message(ERR_NONICKNAMEGIVEN, server.name)
		return
	}

	nickname := args[1]

	if match, _ := regexp.MatchString("^[A-Za-z][\\`_^{|}A-Za-z0-9-]{0,8}$", nickname); !match {
		sender.message(ERR_ERRONEUSNICKNAME, server.name, nickname)
	}

	if server.client(nickname) != nil {
		sender.message(ERR_NICKNAMEINUSE, server.name, nickname)
		return
	}

	sender.nickname = nickname

	if sender.username != "" {
		sender.registered = true
	}
}

func (sender *client) user(server *server, args ...string) {
	args = strings.SplitN(args[1], " ", 4)

	if len(args) < 4 {
		sender.message(ERR_NEEDMOREPARAMS, server.name, sender.nickname, "USER")
		return
	}

	if sender.registered {
		sender.message(ERR_ALREADYREGISTRED, server.name, sender.nickname)
		return
	}

	sender.username = args[0]
	sender.realname = args[3]

	if sender.nickname != "" {
		sender.registered = true
	}
}

func (sender *client) oper(server *server, args ...string) {
	args = strings.SplitN(args[1], " ", 2)

	if len(args) < 2 {
		sender.message(ERR_NEEDMOREPARAMS, server.name, sender.nickname, "OPER")
		return
	}

	if server.login(sender, args[0], args[1]) {
		sender.message(RPL_YOUREOPER, server.name)
	} else {
		sender.message(ERR_PASSWDMISMATCH, server.name)
	}
}

func (sender *client) privmsg(server *server, args ...string) {
	// :martijn PRIVMSG amber :Hi there
	message := strings.SplitN(args[1], " ", 2)
	verbose := args[0] == "PRIVMSG"

	if len(args) < 1 && verbose {
		sender.message(ERR_NORECIPIENT, server.name, sender.nickname, args[0])
		return
	}

	if len(args) < 2 && verbose {
		sender.message(ERR_NOTEXTTOSEND, server.name, sender.nickname)
		return
	}

	// TODO split multiple recipients; martijn,amber
	recipient := message[0]
	text := message[1]

	if channel := server.channel(recipient); channel != nil {
		channel.broadcast(server, sender, ":%s %s %s :%s", sender.prefix(server), args[0], channel.name, text)
		return
	}

	if client := server.client(recipient); client != nil {
		client.message(":%s %s %s :%s", sender.prefix(server), args[0], client.nickname, text)
		return
	}

	if verbose {
		sender.message(ERR_NOSUCHNICK, server.name, sender.nickname, recipient)
	}
}

func (sender *client) join(server *server, args ...string) {
	args = strings.SplitN(args[1], " ", 2)

	if len(args) < 1 {
		sender.message(ERR_NEEDMOREPARAMS, server.name, sender.nickname, "JOIN")
		return
	}

	for _, name := range strings.Split(args[0], ",") {
		name = strings.ToLower(name)

		if match, _ := regexp.MatchString("^[#$][^\x00\x07\x0a\x0d ,:]{1,40}$", name); !match {
			sender.message(ERR_NOSUCHCHANNEL, server.name, sender.nickname, name)
			continue
		}

		channel := server.channels[name]

		if channel == nil {
			sender.message(ERR_NOSUCHCHANNEL, server.name, sender.nickname, name)
			continue
		}

		channel.join(server, sender)
	}
}

func (sender *client) part(server *server, args ...string) {
	args = strings.SplitN(args[1], " ", 2)

	if len(args) < 1 {
		sender.message(ERR_NEEDMOREPARAMS, server.name, sender.nickname, "PART")
		return
	}

	for _, name := range strings.Split(args[0], ",") {
		channel := server.channels[name]

		if channel == nil {
			sender.message(ERR_NOSUCHCHANNEL, server.name, sender.nickname, name)
			continue
		}

		channel.part(server, sender)
	}
}
