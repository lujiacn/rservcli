package constants

type Command int

const (
	CmdLogin    Command = 1
	CmdVoidEval Command = 2
	CmdEval     Command = 3
	CmdSetSexp  Command = 32
)
