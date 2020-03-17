package connection

import (
	"log"
	"strconv"

	"github.com/fanliao/go-promise"

	"main/connection/core"
)

type PetInfoNature uint32

func (n PetInfoNature) String() string {
	switch n {
	case 0:
		return "孤独"
	case 1:
		return "固执"
	case 2:
		return "调皮"
	case 3:
		return "勇敢"
	case 4:
		return "大胆"
	case 5:
		return "顽皮"
	case 6:
		return "无虑"
	case 7:
		return "悠闲"
	case 8:
		return "保守"
	case 9:
		return "稳重"
	case 10:
		return "马虎"
	case 11:
		return "冷静"
	case 12:
		return "沉着"
	case 13:
		return "温顺"
	case 14:
		return "慎重"
	case 15:
		return "狂妄"
	case 16:
		return "胆小"
	case 17:
		return "急躁"
	case 18:
		return "开朗"
	case 19:
		return "天真"
	case 20:
		return "害羞"
	case 21:
		return "实干"
	case 22:
		return "坦率"
	case 23:
		return "浮躁"
	case 24:
		return "认真"
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
			var len uint32
			core.MustBinaryRead(buffer, &len)
			petList := make([]PetListInfo, len)
			for i := uint32(0); i < len; i++ {
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
		info.Name = string(name[:16])
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
