package command

import (
	"fmt"
	"strconv"
	"strings"
)

type PartCommand struct {
	BaseCommand
	channels []Name
	message  Text
}

func (cmd *PartCommand) Message() Text {
	if cmd.message == "" {
		return cmd.Client().Nick().Text()
	}
	return cmd.message
}

func ParsePartCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	msg := &PartCommand{
		channels: NewNames(strings.Split(args[0], ",")),
	}
	if len(args) > 1 {
		msg.message = NewText(args[1])
	}
	return msg, nil
}

// PRIVMSG <target> <message>

type PrivMsgCommand struct {
	BaseCommand
	target  Name
	message Text
}

func ParsePrivMsgCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}
	return &PrivMsgCommand{
		target:  NewName(args[0]),
		message: NewText(args[1]),
	}, nil
}

// TOPIC [newtopic]

type TopicCommand struct {
	BaseCommand
	channel  Name
	setTopic bool
	topic    Text
}

func ParseTopicCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	msg := &TopicCommand{
		channel: NewName(args[0]),
	}
	if len(args) > 1 {
		msg.setTopic = true
		msg.topic = NewText(args[1])
	}
	return msg, nil
}

type ModeChange struct {
	mode UserMode
	op   ModeOp
}

func (change *ModeChange) String() string {
	return fmt.Sprintf("%s%s", change.op, change.mode)
}

type ModeChanges []*ModeChange

func (changes ModeChanges) String() string {
	if len(changes) == 0 {
		return ""
	}

	op := changes[0].op
	str := changes[0].op.String()
	for _, change := range changes {
		if change.op == op {
			str += change.mode.String()
		} else {
			op = change.op
			str += " " + change.op.String()
		}
	}
	return str
}

type ModeCommand struct {
	BaseCommand
	nickname Name
	changes  ModeChanges
}

// MODE <nickname> *( ( "+" / "-" ) *( "i" / "w" / "o" / "O" / "r" ) )
func ParseUserModeCommand(nickname Name, args []string) (Command, error) {
	cmd := &ModeCommand{
		nickname: nickname,
		changes:  make(ModeChanges, 0),
	}

	for _, modeChange := range args {
		if len(modeChange) == 0 {
			continue
		}
		op := ModeOp(modeChange[0])
		if (op != Add) && (op != Remove) {
			return nil, ErrParseCommand
		}

		for _, mode := range modeChange[1:] {
			cmd.changes = append(cmd.changes, &ModeChange{
				mode: UserMode(mode),
				op:   op,
			})
		}
	}

	return cmd, nil
}

type ChannelModeChange struct {
	mode ChannelMode
	op   ModeOp
	arg  string
}

func (change *ChannelModeChange) String() (str string) {
	if (change.op == Add) || (change.op == Remove) {
		str = change.op.String()
	}
	str += change.mode.String()
	if change.arg != "" {
		str += " " + change.arg
	}
	return
}

type ChannelModeChanges []*ChannelModeChange

func (changes ChannelModeChanges) String() (str string) {
	if len(changes) == 0 {
		return
	}

	str = "+"
	if changes[0].op == Remove {
		str = "-"
	}
	for _, change := range changes {
		str += change.mode.String()
	}
	for _, change := range changes {
		if change.arg == "" {
			continue
		}
		str += " " + change.arg
	}
	return
}

type ChannelModeCommand struct {
	BaseCommand
	channel Name
	changes ChannelModeChanges
}

// MODE <channel> *( ( "-" / "+" ) *<modes> *<modeparams> )
func ParseChannelModeCommand(channel Name, args []string) (Command, error) {
	cmd := &ChannelModeCommand{
		channel: channel,
		changes: make(ChannelModeChanges, 0),
	}

	for len(args) > 0 {
		if len(args[0]) == 0 {
			args = args[1:]
			continue
		}

		modeArg := args[0]
		op := ModeOp(modeArg[0])
		if (op == Add) || (op == Remove) {
			modeArg = modeArg[1:]
		} else {
			op = List
		}

		skipArgs := 1
		for _, mode := range modeArg {
			change := &ChannelModeChange{
				mode: ChannelMode(mode),
				op:   op,
			}
			switch change.mode {
			case Key, BanMask, ExceptMask, InviteMask, UserLimit,
				ChannelOperator, ChannelCreator, Voice:
				if len(args) > skipArgs {
					change.arg = args[skipArgs]
					skipArgs += 1
				}
			}
			cmd.changes = append(cmd.changes, change)
		}
		args = args[skipArgs:]
	}

	return cmd, nil
}

func ParseModeCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, NotEnoughArgsError
	}

	name := NewName(args[0])
	if name.IsChannel() {
		return ParseChannelModeCommand(name, args[1:])
	} else {
		return ParseUserModeCommand(name, args[1:])
	}
}

type WhoisCommand struct {
	BaseCommand
	target Name
	masks  []Name
}

// WHOIS [ <target> ] <mask> *( "," <mask> )
func ParseWhoisCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}

	var masks string
	var target string

	if len(args) > 1 {
		target = args[0]
		masks = args[1]
	} else {
		masks = args[0]
	}

	return &WhoisCommand{
		target: NewName(target),
		masks:  NewNames(strings.Split(masks, ",")),
	}, nil
}

type WhoCommand struct {
	BaseCommand
	mask         Name
	operatorOnly bool
}

// WHO [ <mask> [ "o" ] ]
func ParseWhoCommand(args []string) (Command, error) {
	cmd := &WhoCommand{}

	if len(args) > 0 {
		cmd.mask = NewName(args[0])
	}

	if (len(args) > 1) && (args[1] == "o") {
		cmd.operatorOnly = true
	}

	return cmd, nil
}

type OperCommand struct {
	PassCommand
	name Name
}

func (msg *OperCommand) LoadPassword(server *Server) {
	msg.hash = server.operators[msg.name]
}

// OPER <name> <password>
func ParseOperCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}

	cmd := &OperCommand{
		name: NewName(args[0]),
	}
	cmd.password = []byte(args[1])
	return cmd, nil
}

type RehashCommand struct {
	BaseCommand
}

// REHASH
func ParseRehashCommand(args []string) (Command, error) {
	return &RehashCommand{}, nil
}

type CapCommand struct {
	BaseCommand
	subCommand   CapSubCommand
	capabilities CapabilitySet
}

func ParseCapCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}

	cmd := &CapCommand{
		subCommand:   CapSubCommand(strings.ToUpper(args[0])),
		capabilities: make(CapabilitySet),
	}

	if len(args) > 1 {
		strs := spacesExpr.Split(args[1], -1)
		for _, str := range strs {
			cmd.capabilities[Capability(str)] = true
		}
	}
	return cmd, nil
}

type AwayCommand struct {
	BaseCommand
	text Text
}

func ParseAwayCommand(args []string) (Command, error) {
	cmd := &AwayCommand{}

	if len(args) > 0 {
		cmd.text = NewText(args[0])
	}

	return cmd, nil
}

type IsOnCommand struct {
	BaseCommand
	nicks []Name
}

func ParseIsOnCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, NotEnoughArgsError
	}

	return &IsOnCommand{
		nicks: NewNames(args),
	}, nil
}

type MOTDCommand struct {
	BaseCommand
	target Name
}

func ParseMOTDCommand(args []string) (Command, error) {
	cmd := &MOTDCommand{}
	if len(args) > 0 {
		cmd.target = NewName(args[0])
	}
	return cmd, nil
}

type NoticeCommand struct {
	BaseCommand
	target  Name
	message Text
}

func ParseNoticeCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}
	return &NoticeCommand{
		target:  NewName(args[0]),
		message: NewText(args[1]),
	}, nil
}

type KickCommand struct {
	BaseCommand
	kicks   map[Name]Name
	comment Text
}

func (msg *KickCommand) Comment() Text {
	if msg.comment == "" {
		return msg.Client().Nick().Text()
	}
	return msg.comment
}

func ParseKickCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}
	channels := NewNames(strings.Split(args[0], ","))
	users := NewNames(strings.Split(args[1], ","))
	if (len(channels) != len(users)) && (len(users) != 1) {
		return nil, NotEnoughArgsError
	}
	cmd := &KickCommand{
		kicks: make(map[Name]Name),
	}
	for index, channel := range channels {
		if len(users) == 1 {
			cmd.kicks[channel] = users[0]
		} else {
			cmd.kicks[channel] = users[index]
		}
	}
	if len(args) > 2 {
		cmd.comment = NewText(args[2])
	}
	return cmd, nil
}

type ListCommand struct {
	BaseCommand
	channels []Name
	target   Name
}

func ParseListCommand(args []string) (Command, error) {
	cmd := &ListCommand{}
	if len(args) > 0 {
		cmd.channels = NewNames(strings.Split(args[0], ","))
	}
	if len(args) > 1 {
		cmd.target = NewName(args[1])
	}
	return cmd, nil
}

type NamesCommand struct {
	BaseCommand
	channels []Name
	target   Name
}

func ParseNamesCommand(args []string) (Command, error) {
	cmd := &NamesCommand{}
	if len(args) > 0 {
		cmd.channels = NewNames(strings.Split(args[0], ","))
	}
	if len(args) > 1 {
		cmd.target = NewName(args[1])
	}
	return cmd, nil
}

type VersionCommand struct {
	BaseCommand
	target Name
}

func ParseVersionCommand(args []string) (Command, error) {
	cmd := &VersionCommand{}
	if len(args) > 0 {
		cmd.target = NewName(args[0])
	}
	return cmd, nil
}

type InviteCommand struct {
	BaseCommand
	nickname Name
	channel  Name
}

func ParseInviteCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}

	return &InviteCommand{
		nickname: NewName(args[0]),
		channel:  NewName(args[1]),
	}, nil
}

type TimeCommand struct {
	BaseCommand
	target Name
}

func ParseTimeCommand(args []string) (Command, error) {
	cmd := &TimeCommand{}
	if len(args) > 0 {
		cmd.target = NewName(args[0])
	}
	return cmd, nil
}

type LUsersCommand struct {
	BaseCommand
}

func ParseLUsersCommand(args []string) (Command, error) {
	return &LUsersCommand{}, nil
}

type KillCommand struct {
	BaseCommand
	nickname Name
	comment  Text
}

func ParseKillCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}
	return &KillCommand{
		nickname: NewName(args[0]),
		comment:  NewText(args[1]),
	}, nil
}

type WallopsCommand struct {
	BaseCommand
	message Text
}

func ParseWallopsCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	return &WallopsCommand{
		message: NewText(args[0]),
	}, nil
}

type WhoWasCommand struct {
	BaseCommand
	nicknames []Name
	count     int64
	target    Name
}

func ParseWhoWasCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	cmd := &WhoWasCommand{
		nicknames: NewNames(strings.Split(args[0], ",")),
	}
	if len(args) > 1 {
		cmd.count, _ = strconv.ParseInt(args[1], 10, 64)
	}
	if len(args) > 2 {
		cmd.target = NewName(args[2])
	}
	return cmd, nil
}

func ParseOperNickCommand(args []string) (Command, error) {
	if len(args) < 2 {
		return nil, NotEnoughArgsError
	}

	return &OperNickCommand{
		target: NewName(args[0]),
		nick:   NewName(args[1]),
	}, nil
}
