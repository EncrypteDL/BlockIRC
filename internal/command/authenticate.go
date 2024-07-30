package command

type AuthenticateCommand struct {
	BaseCommand
	arg string
}

func ParseAuthenticateCommand(args []string) (Command, error) {
	if len(args) < 1 {
		return nil, NotEnoughArgsError
	}
	return &AuthenticateCommand{
		arg: args[0],
	}, nil
}