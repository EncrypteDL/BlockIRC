package command

type PassCommand struct {
	BaseCommand
	hash     []byte
	password []byte
	err      error
}

func (cmd *PassCommand) LoadPassword(server *Server) {
	cmd.hash = server.password
}

func (cmd *PassCommand) CheckPassword() {
	if cmd.hash == nil {
		return
	}
	cmd.err = ComparePassword(cmd.hash, cmd.password)
}

func ParsePassCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	return &PassCommand{
		password: []byte(args[0]),
	}, nil
}
