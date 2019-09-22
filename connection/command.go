package connection

import "main/connection/core"

//noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	Command_GET_VERIFCODE  core.Command = 101
	Command_MAIN_LOGIN_IN  core.Command = 103
	Command_COMMEND_ONLINE core.Command = 105
	Command_RANGE_ONLINE   core.Command = 106
	Command_CREATE_ROLE    core.Command = 108
	Command_SYS_ROLE       core.Command = 109

	Command_LOGIN_IN      core.Command = 1001
	Command_COMPLETE_TASK core.Command = 2202
	Command_Test          core.Command = 30000
)
