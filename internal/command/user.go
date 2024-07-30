package command

import (
	"strconv"
	"strings"
)

func ParseNickCommand(args []string) (Command, error) {
	if len(args) != 1 {
		return nil, NotEnoughArgsError
	}
	return &NickCommand{
		nickname: NewName(args[0]),
	}, nil
}

type UserCommand struct {
	BaseCommand
	username Name
	realname Text
}

// USER <username> <hostname> <servername> <realname>
type RFC1459UserCommand struct {
	UserCommand
	hostname   Name
	servername Name
}

// USER <user> <mode> <unused> <realname>
type RFC2812UserCommand struct {
	UserCommand
	mode   uint8
	unused string
}

func (cmd *RFC2812UserCommand) Flags() []UserMode {
	flags := make([]UserMode, 0)
	if (cmd.mode & 4) == 4 {
		flags = append(flags, WallOps)
	}
	if (cmd.mode & 8) == 8 {
		flags = append(flags, Invisible)
	}
	return flags
}

func ParseUserCommand(args []string) (Command, error) {
	if len(args) != 4 {
		return nil, NotEnoughArgsError
	}
	mode, err := strconv.ParseUint(args[1], 10, 8)
	if err == nil {
		msg := &RFC2812UserCommand{
			mode:   uint8(mode),
			unused: args[2],
		}
		msg.username = NewName(args[0])
		msg.realname = NewText(args[3])
		return msg, nil
	}

	msg := &RFC1459UserCommand{
		hostname:   NewName(args[1]),
		servername: NewName(args[2]),
	}
	msg.username = NewName(args[0])
	msg.realname = NewText(args[3])
	return msg, nil
}

// QUIT [ <Quit Command> ]

type QuitCommand struct {
	BaseCommand
	message Text
}

func NewQuitCommand(message Text) *QuitCommand {
	cmd := &QuitCommand{
		message: message,
	}
	cmd.code = QUIT
	return cmd
}

func ParseQuitCommand(args []string) (Command, error) {
	msg := &QuitCommand{}
	if len(args) > 0 {
		msg.message = NewText(args[0])
	}
	return msg, nil
}

// JOIN ( <channel> *( "," <channel> ) [ <key> *( "," <key> ) ] ) / "0"

type JoinCommand struct {
	BaseCommand
	channels map[Name]Text
	zero     bool
}

func ParseJoinCommand(args []string) (Command, error) {
	msg := &JoinCommand{
		channels: make(map[Name]Text),
	}

	if len(args) == 0 {
		return nil, NotEnoughArgsError
	}

	if args[0] == "0" {
		msg.zero = true
		return msg, nil
	}

	channels := strings.Split(args[0], ",")
	keys := make([]string, len(channels))
	if len(args) > 1 {
		for i, key := range strings.Split(args[1], ",") {
			if i >= len(channels) {
				break
			}
			keys[i] = key
		}
	}
	for i, channel := range channels {
		msg.channels[NewName(channel)] = NewText(keys[i])
	}

	return msg, nil
}
