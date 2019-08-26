package connection

type Command = uint32

const (
	Command_GET_VERIFCODE  Command = 101
	Command_MAIN_LOGIN_IN  Command = 103
	Command_COMMEND_ONLINE Command = 105
	Command_RANGE_ONLINE   Command = 106
	Command_CREATE_ROLE    Command = 108
	Command_SYS_ROLE       Command = 109

	Command_LOGIN_IN Command = 1001
	Command_Test     Command = 30000
)
