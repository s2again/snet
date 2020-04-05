package snet

import (
	"github.com/fanliao/go-promise"
)

type RoleColor = uint32

//noinspection GoUnusedConst
const (
	RoleYellow RoleColor = 0xFFFF00
	RoleGreen  RoleColor = 0x00FF00
	RoleRed    RoleColor = 0xFD3E3E
	RoleOrange RoleColor = 0xFF6500
	RoleBlue   RoleColor = 0x0000FF
	RoleBrown  RoleColor = 0x996600
	RoleWhite  RoleColor = 0xFFFFFF
	RoleBlack  RoleColor = 0x000000
)

func (c *GuideServerConnection) CreateRole(nickname [16]byte, color RoleColor) *promise.Promise {
	const verifyCode uint32 = 0 // 邀请码机制，已作废。传0即可。
	return c.SendInPromise(Command_CREATE_ROLE, c.SessionID, nickname, color, verifyCode, CHANNEL)
}
