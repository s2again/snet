package connection

import (
	"github.com/fanliao/go-promise"

	"main/connection/core"
)

type NoviceFinishItem struct {
	ItemID  uint32
	ItemCnt uint32
}
type NoviceFinishInfo struct {
	TaskID       uint32
	PetID        uint32
	CaptureTm    uint32 // capture timestamp
	itemListSize uint32
	ItemList     []NoviceFinishItem
}

func (c *Connection) AcceptTask(taskID uint32) *promise.Promise {
	return c.SendInPromise(Command_ACCEPT_TASK, taskID)
}

func (c *Connection) CompleteTask(taskID uint32, param uint32) *promise.Promise {
	prom := promise.NewPromise()
	c.SendInPromise(Command_COMPLETE_TASK, taskID, param).
		OnSuccess(func(v interface{}) {
			info, err := parseNoviceFinishInfo(v.(core.PacketBody))
			if err != nil {
				prom.Reject(v.(error))
			} else {
				prom.Resolve(info)
			}
		}).
		OnFailure(func(v interface{}) {
			prom.Reject(v.(error))
		})
	return prom
}

func parseNoviceFinishInfo(buffer core.PacketBody) (info NoviceFinishInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	core.MustBinaryRead(buffer, &info.TaskID)
	core.MustBinaryRead(buffer, &info.PetID)
	core.MustBinaryRead(buffer, &info.CaptureTm)
	core.MustBinaryRead(buffer, &info.itemListSize)
	info.ItemList = make([]NoviceFinishItem, info.itemListSize)
	for i := uint32(0); i < info.itemListSize; i++ {
		core.MustBinaryRead(buffer, &info.ItemList[i].ItemID)
		core.MustBinaryRead(buffer, &info.ItemList[i].ItemCnt)
	}
	return info, nil
}
