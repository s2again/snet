package snet

import "github.com/s2again/snet/core"

//noinspection GoUnusedConst,GoSnakeCaseUsage
const (
	Command_SEER_VERIFY core.Command = 1
	Command_REGISTER    core.Command = 2

	// Command_GET_VERIFCODE core.Command = 101
	// Command_MAIN_LOGIN_IN  core.Command = 103
	Command_MAIN_LOGIN_IN  core.Command = 104
	Command_COMMEND_ONLINE core.Command = 105 // 0x69
	Command_RANGE_ONLINE   core.Command = 106
	Command_CREATE_ROLE    core.Command = 108 // 0x6c
	// Command_SYS_ROLE       core.Command = 109

	Command_LOGIN_IN    core.Command = 1001 // 0x03E9
	Command_SYSTEM_TIME core.Command = 1002 // 0x03EA

	Command_GET_SESSION_KEY core.Command = 1006

	Command_ENTER_MAP       core.Command = 2001 // 0x07D1
	Command_LEAVE_MAP       core.Command = 2002 // 0x07D2
	Command_LIST_MAP_PLAYER core.Command = 2003 // 0x07D3

	Command_ACCEPT_TASK   core.Command = 2201 // 0x0899
	Command_COMPLETE_TASK core.Command = 2202 // 0x089A
	// ...

	Command_GET_PET_INFO    core.Command = 2301 // 0x08FD
	Command_MODIFY_PET_NAME core.Command = 2302
	Command_GET_PET_LIST    core.Command = 2303
	Command_PET_RELEASE     core.Command = 2304 // 0x0900
	// ...

	Command_GET_SOUL_BEAD_LIST core.Command = 2354 // 0x0932

	Command_MAIL_GET_LIST core.Command = 2751
	// ...
	Command_MAIL_GET_UNREAD core.Command = 2757 // 0x0AC5
	// ...

	Command_NONO_OPEN        core.Command = 9001
	Command_NONO_CHANGE_NAME core.Command = 9002
	Command_NONO_INFO        core.Command = 9003 // 0x232B
	// ...

	Command_TEST core.Command = 30000
)

const CHANNEL uint32 = 0 // MainManager.CHANNEL
