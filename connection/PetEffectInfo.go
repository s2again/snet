package connection

import "main/connection/core"

// com.robot.core.info.pet.PetEffectInfo
type PetEffectInfo struct {
	ItemID    int32
	Status    byte
	LeftCount byte
	EffectID  uint16
	Arg       [16]byte
}

// com.robot.core.info.pet.PetEffectInfo
func parsePetEffectInfo(buffer core.PacketBody) (info PetEffectInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	core.MustBinaryRead(buffer, &info.ItemID)
	core.MustBinaryRead(buffer, &info.Status)
	core.MustBinaryRead(buffer, &info.LeftCount)
	core.MustBinaryRead(buffer, &info.EffectID)
	core.MustBinaryRead(buffer, &info.Arg)
	return
}
