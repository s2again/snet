package connection

import (
	"github.com/fanliao/go-promise"
)

type RoleColor = uint32

var (
	RoleGreen RoleColor = 0x0000FF00
)

func (c *Connection) CreateRole(nickname [16]byte, color RoleColor) *promise.Promise {
	const verifyCode uint32 = 0 // 邀请码机制，已作废。传0即可。
	return c.SendInPromise(Command_CREATE_ROLE, c.SessionID, nickname, color, verifyCode, CHANNEL)
}
