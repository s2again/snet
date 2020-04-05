package snet

import (
	"bytes"
	"log"
	"time"

	"github.com/fanliao/go-promise"

	"main/snet/core"
)

type ResponseForLogin struct {
	UserID                                                        uint32
	RegTime                                                       time.Time
	Nick                                                          string
	Vip                                                           bool
	Viped                                                         bool
	Color, Texture, Energy, Coins, FightBadge                     uint32
	MapID, PosX, PosY                                             uint32
	TimeToday, TimeLimit                                          uint32
	IsClothHalfDay, IsRoomHalfDay, IsFortressHalfDay, IsHQHalfDay bool
	LoginCnt                                                      uint32
}

func (c *OnlineServerConnection) LoginOnline() *promise.Promise {
	prom := promise.NewPromise()
	c.SendInPromise(Command_LOGIN_IN, c.SessionID).OnSuccess(func(v interface{}) {
		info, err := parseUserInfoForLogin(v.(core.PacketBody))
		if err != nil {
			log.Println("parseCommendSvrInfo error: ", err)
			prom.Reject(err)
			return
		}
		prom.Resolve(info)
	}).OnFailure(func(v interface{}) {
		log.Println("LoginOnlineAndCallback promise rejected: ", v)
		prom.Reject(v.(error))
	})
	return prom
}

func (c *OnlineServerConnection) LoginOnlineAndCallback(callback func(ResponseForLogin)) error {
	body, err := c.SendInPromise(Command_LOGIN_IN, c.SessionID).Get()
	if err != nil {
		return err
	}
	info, err := parseUserInfoForLogin(body.(core.PacketBody))
	if err != nil {
		log.Println("parseCommendSvrInfo error: ", err, "connection terminated.")
		return err
	}
	callback(info)
	return nil
}

func parseUserInfoForLogin(buffer core.PacketBody) (info ResponseForLogin, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	var (
		userID, regTime                                              uint32
		nick                                                         [16]byte
		vipInfo, dsFlag                                              uint32
		color, texture, energy, coins, fightBadge                    uint32
		mapID, posX, posY                                            uint32
		timeToday, timeLimit                                         uint32
		isClothHalfDay, isRoomHalfDay, iFortressHalfDay, isHQHalfDay byte
		loginCnt                                                     uint32
		inviter, newInviteeCnt                                       uint32
		vipLevel, vipValue, vipStage, autoCharge, vipEndTime         uint32
		freshManBonus                                                uint32
		nonoChipList                                                 [80]byte
		dailyRes                                                     [50]byte
		// ...
	)
	core.MustBinaryRead(buffer, &userID, &regTime, &nick)
	core.MustBinaryRead(buffer, &vipInfo, &dsFlag)
	core.MustBinaryRead(buffer, &color, &texture, &energy, &coins, &fightBadge)
	core.MustBinaryRead(buffer, &mapID, &posX, &posY)
	core.MustBinaryRead(buffer, &timeToday, &timeLimit)
	core.MustBinaryRead(buffer, &isClothHalfDay, &isRoomHalfDay, &iFortressHalfDay, &isHQHalfDay)
	core.MustBinaryRead(buffer, &loginCnt, &inviter, &newInviteeCnt)
	core.MustBinaryRead(buffer, &vipLevel, &vipValue, &vipStage, &autoCharge, &vipEndTime)
	core.MustBinaryRead(buffer, &freshManBonus, &nonoChipList, &dailyRes)
	return ResponseForLogin{
		UserID:  userID,
		RegTime: time.Unix(int64(regTime), 0),
		Nick:    string(bytes.Trim(nick[:], "\u0000")),
		// TODO
		// Vip = vipInfo >> ?
		// Viped = vipInfo >> ?
		Color: color, Texture: texture, Energy: energy, Coins: coins, FightBadge: fightBadge,
		MapID: mapID, PosX: posX, PosY: posY,
		TimeToday: timeToday, TimeLimit: timeLimit,
		IsClothHalfDay:    isClothHalfDay != 0,
		IsRoomHalfDay:     isRoomHalfDay != 0,
		IsFortressHalfDay: iFortressHalfDay != 0,
		IsHQHalfDay:       isHQHalfDay != 0,
		LoginCnt:          loginCnt,
	}, nil
	/*
		参考解析代码：
		        var id:uint = 0;
		         var level:uint = 0;
		         info.hasSimpleInfo = true;
		         info.userID = data.readUnsignedInt();
		         info.regTime = data.readUnsignedInt();
		         info.nick = data.readUTFBytes(16);
		         var vvv:uint = data.readUnsignedInt();
		         info.vip = BitUtil.getBit(vvv,0);
		         info.viped = BitUtil.getBit(vvv,1);
		         info.dsFlag = data.readUnsignedInt();
		         info.color = data.readUnsignedInt();
		         info.texture = data.readUnsignedInt();
		         info.energy = data.readUnsignedInt();
		         info.coins = data.readUnsignedInt();
		         info.fightBadge = data.readUnsignedInt();
		         info.mapID = data.readUnsignedInt();
		         info.pos = new Point(data.readUnsignedInt(),data.readUnsignedInt());
		         info.timeToday = data.readUnsignedInt();
		         info.timeLimit = data.readUnsignedInt();
		         MainManager.isClothHalfDay = Boolean(data.readByte());
		         MainManager.isRoomHalfDay = Boolean(data.readByte());
		         MainManager.iFortressHalfDay = Boolean(data.readByte());
		         MainManager.isHQHalfDay = Boolean(data.readByte());
		         trace("个人装扮是否半价：",MainManager.isClothHalfDay);
		         trace("小屋装扮是否半价：",MainManager.isRoomHalfDay);
		         trace("要塞装扮是否半价：",MainManager.iFortressHalfDay);
		         trace("总部装扮是否半价：",MainManager.isHQHalfDay);
		         info.loginCnt = data.readUnsignedInt();
		         info.inviter = data.readUnsignedInt();
		         info.newInviteeCnt = data.readUnsignedInt();
		         info.vipLevel = data.readUnsignedInt();
		         info.vipValue = data.readUnsignedInt();
		         info.vipStage = data.readUnsignedInt();
		         if(info.vipStage > 4)
		         {
		            info.vipStage = 4;
		         }
		         if(info.vipStage == 0)
		         {
		            info.vipStage = 1;
		         }
		         info.autoCharge = data.readUnsignedInt();
		         info.vipEndTime = data.readUnsignedInt();
		         info.freshManBonus = data.readUnsignedInt();
		         for(var r:int = 0; r < 80; r++)
		         {
		            info.nonoChipList.push(Boolean(data.readByte()));
		         }
		         for(var rr:int = 0; rr < 50; rr++)
		         {
		            info.dailyResArr.push(data.readByte());
		         }
		         info.teacherID = data.readUnsignedInt();
		         info.studentID = data.readUnsignedInt();
		         info.graduationCount = data.readUnsignedInt();
		         info.maxPuniLv = data.readUnsignedInt();
		         info.petMaxLev = data.readUnsignedInt();
		         info.petAllNum = data.readUnsignedInt();
		         info.monKingWin = data.readUnsignedInt();
		         info.curStage = data.readUnsignedInt() + 1;
		         info.maxStage = data.readUnsignedInt();
		         info.curFreshStage = data.readUnsignedInt();
		         info.maxFreshStage = data.readUnsignedInt();
		         info.maxArenaWins = data.readUnsignedInt();
		         info.twoTimes = data.readUnsignedInt();
		         info.threeTimes = data.readUnsignedInt();
		         info.autoFight = data.readUnsignedInt();
		         info.autoFightTimes = data.readUnsignedInt();
		         info.energyTimes = data.readUnsignedInt();
		         info.learnTimes = data.readUnsignedInt();
		         info.monBtlMedal = data.readUnsignedInt();
		         info.recordCnt = data.readUnsignedInt();
		         info.obtainTm = data.readUnsignedInt();
		         info.soulBeadItemID = data.readUnsignedInt();
		         info.expireTm = data.readUnsignedInt();
		         info.fuseTimes = data.readUnsignedInt();
		         info.hasNono = Boolean(data.readUnsignedInt());
		         info.superNono = Boolean(data.readUnsignedInt());
		         var num:uint = data.readUnsignedInt();
		         for(var s:int = 0; s < 32; s++)
		         {
		            info.nonoState.push(BitUtil.getBit(num,s));
		         }
		         info.nonoColor = data.readUnsignedInt();
		         info.nonoNick = data.readUTFBytes(16);
		         info.teamInfo = new TeamInfo(data);
		         info.teamPKInfo = new TeamPKInfo(data);
		         data.readByte();
		         info.badge = data.readUnsignedInt();
		         var byte:ByteArray = new ByteArray();
		         data.readBytes(byte,0,27);
		         info.reserved = byte;
		         for(var i:int = 0; i < 500; i++)
		         {
		            TasksManager.taskList.push(data.readUnsignedByte());
		         }
		         var a:Array = TasksManager.taskList;
		         info.isCanBeTeacher = TasksManager.getTaskStatus(201) == 3;
		         info.petNum = data.readUnsignedInt();
		         PetManager.initData(data,info.petNum);
		         var clothNum:uint = data.readUnsignedInt();
		         for(var j:uint = 0; j < clothNum; j++)
		         {
		            id = data.readUnsignedInt();
		            level = data.readUnsignedInt();
		            info.clothes.push(new PeopleItemInfo(id,level));
		         }
	*/
}
