# simple chat server in go

Started working on an [irc server](https://tools.ietf.org/html/rfc1459) in Go to try out a new language. As there are multiple IRC Go servers I will discontinue this project. I used [miniircd](https://github.com/jrosdahl/miniircd) a lot as a reference. In go Ifound [oragono](https://github.com/oragono/oragono) and [ggircd](https://github.com/fimad/ggircd) which I found hard to read, but Go is new to me. It uses a goroutine per connection. I haven't tested this software at all except with the following commands.

```bash
# start server
git clone git@github.com:martijncasteel/ircd.go.git
cd ircd.go
go run ircd.go

# start tcp connection
telnet localhost 8000
```

```
# register as a user
nick martijn
user martijn - - :Martijn Casteel

# join a channel, or part from it
join #general
part #general

# send messages
privmsg martijn Hallo Martijn
privmsg #general Hallo everybody

# same but without error messages
notice #general Hallo everybody

# close connection
quit
```
