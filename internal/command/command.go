package command

import (
	"errors"
	"regexp"
	"strings"
)

type Command interface {
	Client() *Client
	Code() StringCode
	SetClient(*Client)
	SetCode(StringCode)
}

type checkPasswordCommand interface {
	LoadPassword(*Server)
	CheckPassword()
}

type parseCommandFunc func([]string) (Command, error)

var (
	NotEnoughArgsError = errors.New("not enough arguments")
	ErrParseCommand    = errors.New("failed to parse message")
	parseCommandFuncs  = map[StringCode]parseCommandFunc{
		AUTHENTICATE: ParseAuthenticateCommand,
		AWAY:         ParseAwayCommand,
		CAP:          ParseCapCommand,
		INVITE:       ParseInviteCommand,
		ISON:         ParseIsOnCommand,
		JOIN:         ParseJoinCommand,
		KICK:         ParseKickCommand,
		KILL:         ParseKillCommand,
		LIST:         ParseListCommand,
		MODE:         ParseModeCommand,
		MOTD:         ParseMOTDCommand,
		NAMES:        ParseNamesCommand,
		NICK:         ParseNickCommand,
		NOTICE:       ParseNoticeCommand,
		ONICK:        ParseOperNickCommand,
		OPER:         ParseOperCommand,
		REHASH:       ParseRehashCommand,
		PART:         ParsePartCommand,
		PASS:         ParsePassCommand,
		PING:         ParsePingCommand,
		PONG:         ParsePongCommand,
		PRIVMSG:      ParsePrivMsgCommand,
		QUIT:         ParseQuitCommand,
		TIME:         ParseTimeCommand,
		LUSERS:       ParseLUsersCommand,
		TOPIC:        ParseTopicCommand,
		USER:         ParseUserCommand,
		VERSION:      ParseVersionCommand,
		WALLOPS:      ParseWallopsCommand,
		WHO:          ParseWhoCommand,
		WHOIS:        ParseWhoisCommand,
		WHOWAS:       ParseWhoWasCommand,
	}
)

type BaseCommand struct {
	client *Client
	code   StringCode
}

func (command *BaseCommand) Client() *Client {
	return command.client
}

func (command *BaseCommand) SetClient(client *Client) {
	command.client = client
}

func (command *BaseCommand) Code() StringCode {
	return command.code
}

func (command *BaseCommand) SetCode(code StringCode) {
	command.code = code
}

func ParseCommand(line string) (cmd Command, err error) {
	code, args := ParseLine(line)
	constructor := parseCommandFuncs[code]
	if constructor == nil {
		cmd = ParseUnknownCommand(args)
	} else {
		cmd, err = constructor(args)
	}
	if cmd != nil {
		cmd.SetCode(code)
	}
	return
}

var (
	spacesExpr = regexp.MustCompile(` +`)
)

func splitArg(line string) (arg string, rest string) {
	parts := spacesExpr.Split(line, 2)
	if len(parts) > 0 {
		arg = parts[0]
	}
	if len(parts) > 1 {
		rest = parts[1]
	}
	return
}

func ParseLine(line string) (command StringCode, args []string) {
	args = make([]string, 0)
	if strings.HasPrefix(line, ":") {
		_, line = splitArg(line)
	}
	arg, line := splitArg(line)
	command = StringCode(NewName(strings.ToUpper(arg)))
	for len(line) > 0 {
		if strings.HasPrefix(line, ":") {
			args = append(args, line[len(":"):])
			break
		}
		arg, line = splitArg(line)
		args = append(args, arg)
	}
	return
}

// <command> [args...]

type UnknownCommand struct {
	BaseCommand
	args []string
}

func ParseUnknownCommand(args []string) *UnknownCommand {
	return &UnknownCommand{
		args: args,
	}
}