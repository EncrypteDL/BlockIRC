package command

type PingCommand struct {
	BaseCommand
	server  Name
	server2 Name
}

func ParsePingCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	msg := &PingCommand{
		server: NewName(args[0]),
	}
	if len(args) > 1 {
		msg.server2 = NewName(args[1])
	}
	return msg, nil
}

// PONG <server> [ <server2> ]

type PongCommand struct {
	BaseCommand
	server1 Name
	server2 Name
}

func ParsePongCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	message := &PongCommand{
		server1: NewName(args[0]),
	}
	if len(args) > 1 {
		message.server2 = NewName(args[1])
	}
	return message, nil
}
