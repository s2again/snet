package snet

import "snet/core"

// com.robot.core.info.pet.PetListInfo
type PetListInfo struct {
	ID        uint32
	CatchTime uint32
}

// com.robot.core.manager.PetManager.getStorageList()
func parsePetListInfo(buffer core.PacketBody) (info PetListInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	core.MustBinaryRead(buffer, &info)
	return
}
