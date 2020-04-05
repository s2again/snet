package snet

import "github.com/fanliao/go-promise"

func (c *OnlineServerConnection) EnterMap(mapType, mapID, x, y uint32) *promise.Promise {
	return c.SendInPromise(Command_ENTER_MAP, mapType, mapID, x, y)
}
func (c *OnlineServerConnection) LeaveMap() *promise.Promise {
	return c.SendInPromise(Command_LEAVE_MAP)
}
func (c *OnlineServerConnection) ListMapPlayer() *promise.Promise {
	return c.SendInPromise(Command_LIST_MAP_PLAYER) // TODO: parse the body
}
