package snet

import (
	"log"
	"strconv"

	"github.com/fanliao/go-promise"

	"main/snet/core"
)

type PetInfoNature uint32

func (n PetInfoNature) String() string {
	if n < 25 {
		return []string{"孤独", "固执", "调皮", "勇敢",
			"大胆", "顽皮", "无虑", "悠闲",
			"保守", "稳重", "马虎", "冷静",
			"沉着", "温顺", "慎重", "狂妄",
			"胆小", "急躁", "开朗", "天真",
			"害羞", "实干", "坦率", "浮躁",
			"认真",
		}[n]
	}
	return "性格" + strconv.FormatUint(uint64(n), 10)
}

// com.robot.core.info.pet.PetInfo
type PetInfo struct {
	ID          uint32
	Name        string
	Dv          uint32
	Nature      PetInfoNature
	Level       uint32
	Exp         uint32
	LvExp       uint32
	NextLvExp   uint32
	HP          uint32
	MaxHP       uint32
	Attack      uint32
	Defence     uint32
	SpAttack    uint32
	SpDefence   uint32
	Speed       uint32
	EvHP        uint32
	EvAttack    uint32
	EvDefence   uint32
	EvSpAttack  uint32
	EvSpDefence uint32
	EvSpeed     uint32
	SkillNum    uint32
	SkillList   []PetSkillInfo
	CatchTime   uint32
	CatchMap    uint32
	CatchRect   uint32
	CatchLevel  uint32
	EffectCount uint16
	EffectList  []PetEffectInfo
}

func (c *Connection) ReleasePet(catchTime uint32, flag uint32) *promise.Promise {
	return c.SendInPromise(Command_PET_RELEASE, catchTime, flag) // TODO: parse the body
}

func (c *Connection) GetPetInfo(catchTime uint32) (p *promise.Promise) {
	p = promise.NewPromise()
	c.SendInPromise(Command_GET_PET_INFO, catchTime).
		OnSuccess(func(v interface{}) {
			info, err := parsePetInfo(v.(core.PacketBody))
			if err != nil {
				p.Reject(err)
			} else {
				p.Resolve(info)
			}
		}).
		OnFailure(func(v interface{}) {
			p.Reject(v.(error))
		})
	return p
}

func (c *Connection) GetPetList() (p *promise.Promise) {
	p = promise.NewPromise()
	c.SendInPromise(Command_GET_PET_LIST).
		OnSuccess(func(v interface{}) {
			defer func() {
				if x := recover(); x != nil {
					p.Reject(x.(error))
				}
			}()
			buffer := v.(core.PacketBody)
			var size uint32
			core.MustBinaryRead(buffer, &size)
			petList := make([]PetListInfo, size)
			for i := uint32(0); i < size; i++ {
				var err error
				petList[i], err = parsePetListInfo(buffer)
				if err != nil {
					p.Reject(v.(error))
					return
				}
			}
			p.Resolve(petList)
		}).
		OnFailure(func(v interface{}) {
			p.Reject(v.(error))
		})
	return p
}

// com.robot.core.info.pet.PetInfo
func parsePetInfo(buffer core.PacketBody) (info PetInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	log.Println("parsePetInfo", buffer.Bytes())
	core.MustBinaryRead(buffer, &info.ID)
	{
		var name [16]byte
		core.MustBinaryRead(buffer, &name)
		info.Name = "" // 暂时无法解析数据。（怀疑响应的name是乱码数据）
	}
	core.MustBinaryRead(buffer, &info.Dv)
	core.MustBinaryRead(buffer, &info.Nature)
	core.MustBinaryRead(buffer, &info.Level)
	core.MustBinaryRead(buffer, &info.Exp)
	core.MustBinaryRead(buffer, &info.LvExp)
	core.MustBinaryRead(buffer, &info.NextLvExp)
	core.MustBinaryRead(buffer, &info.HP)
	core.MustBinaryRead(buffer, &info.MaxHP)
	core.MustBinaryRead(buffer, &info.Attack)
	core.MustBinaryRead(buffer, &info.Defence)
	core.MustBinaryRead(buffer, &info.SpAttack)
	core.MustBinaryRead(buffer, &info.SpDefence)
	core.MustBinaryRead(buffer, &info.Speed)
	core.MustBinaryRead(buffer, &info.EvHP)
	core.MustBinaryRead(buffer, &info.EvAttack)
	core.MustBinaryRead(buffer, &info.EvDefence)
	core.MustBinaryRead(buffer, &info.EvSpAttack)
	core.MustBinaryRead(buffer, &info.EvSpDefence)
	core.MustBinaryRead(buffer, &info.EvSpeed)
	core.MustBinaryRead(buffer, &info.SkillNum)
	info.SkillList = make([]PetSkillInfo, 4)
	for i := 0; i < len(info.SkillList); i++ {
		info.SkillList[i], err = parsePetSkillInfo(buffer)
		if err != nil {
			panic(err)
		}
	}
	info.SkillList = info.SkillList[:info.SkillNum]
	core.MustBinaryRead(buffer, &info.CatchTime)
	core.MustBinaryRead(buffer, &info.CatchMap)
	core.MustBinaryRead(buffer, &info.CatchRect)
	core.MustBinaryRead(buffer, &info.Level)
	core.MustBinaryRead(buffer, &info.EffectCount)
	info.EffectList = make([]PetEffectInfo, info.EffectCount)
	for i := uint16(0); i < info.EffectCount; i++ {
		info.EffectList[i], err = parsePetEffectInfo(buffer)
		if err != nil {
			panic(err)
		}
	}
	return info, nil
}
