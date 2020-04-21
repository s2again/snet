package snet

import "github.com/s2again/snet/core"

// com.robot.core.info.pet.PetSkillInfo
type PetSkillInfo struct {
	ID uint32
	PP uint32
}

// com.robot.core.info.pet.PetSkillInfo
func parsePetSkillInfo(buffer core.PacketBody) (info PetSkillInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	core.MustBinaryRead(buffer, &info)
	return
}
