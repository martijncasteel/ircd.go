# simple chat server in golang


``` 
# start tcp connection
> telnet localhost 8000
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

```