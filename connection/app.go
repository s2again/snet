package connection

import (
	"bytes"
	"log"
	"time"
)

const gameChannel uint32 = 0

// 当前不校验session有效性，因此调用者自行保证其有效性。
func (c *Connection) SetSession(userID uint32, sessionID [16]byte) {
	c.UserID, c.SessionID = userID, sessionID
}

func (c *Connection) ListOnlineServers(getCommendList func(CommendSvrInfo)) error {
	var id MsgListenerID
	id = c.AddListener(Command_COMMEND_ONLINE, func(body packetBody) {
		c.RemoveListener(Command_COMMEND_ONLINE, id)
		info, err := parseCommendSvrInfo(body)
		if err != nil {
			_ = c.Close()
		}
		getCommendList(info)
	})
	err := c.Send(Command_COMMEND_ONLINE, c.SessionID, gameChannel)
	if err != nil {
		return err
	}
	return nil
}

type OnlineServerInfo struct {
	OnlineID uint32
	UserCnt  uint32
	IP       string
	Port     uint16
	Friends  uint32
}
type CommendSvrInfo struct {
	MaxOnlineID uint32
	IsVIP       uint32
	OnlineCnt   uint32
	SvrList     []OnlineServerInfo
	// friendList []byte
}

func parseCommendSvrInfo(buffer packetBody) (info CommendSvrInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	log.Println("Command_COMMEND_ONLINE response bytes", buffer.Bytes())
	mustBinaryRead(buffer, &info.MaxOnlineID)
	mustBinaryRead(buffer, &info.IsVIP)
	mustBinaryRead(buffer, &info.OnlineCnt)
	log.Println("OnlineCnt", info.OnlineCnt)
	info.SvrList = make([]OnlineServerInfo, info.OnlineCnt)
	for i := uint32(0); i < info.OnlineCnt; i++ {
		mustBinaryRead(buffer, &info.SvrList[i].OnlineID)
		mustBinaryRead(buffer, &info.SvrList[i].UserCnt)
		{
			var ipBin [16]byte
			mustBinaryRead(buffer, &ipBin)
			info.SvrList[i].IP = string(bytes.Trim(ipBin[:], "\u0000"))
		}
		mustBinaryRead(buffer, &info.SvrList[i].Port)
		mustBinaryRead(buffer, &info.SvrList[i].Friends)
	}
	return
}

type UserInfo struct {
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

func (c *Connection) LoginOnline(userInfoFunc func(UserInfo)) error {
	var id MsgListenerID
	id = c.AddListener(Command_LOGIN_IN, func(body packetBody) {
		c.RemoveListener(Command_LOGIN_IN, id)
		log.Println("LoginOnline resp", body.Bytes())
		info, err := parseUserInfoForLogin(body)
		if err != nil {
			_ = c.Close()
		}
		userInfoFunc(info)
	})
	err := c.Send(Command_LOGIN_IN, c.SessionID)
	if err != nil {
		return err
	}
	return nil
}

func parseUserInfoForLogin(buffer packetBody) (info UserInfo, err error) {
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
	mustBinaryRead(buffer, &userID, &regTime, &nick)
	mustBinaryRead(buffer, &vipInfo, &dsFlag)
	mustBinaryRead(buffer, &color, &texture, &energy, &coins, &fightBadge)
	mustBinaryRead(buffer, &mapID, &posX, &posY)
	mustBinaryRead(buffer, &timeToday, &timeLimit)
	mustBinaryRead(buffer, &isClothHalfDay, &isRoomHalfDay, &iFortressHalfDay, &isHQHalfDay)
	mustBinaryRead(buffer, &loginCnt, &inviter, &newInviteeCnt)
	mustBinaryRead(buffer, &vipLevel, &vipValue, &vipStage, &autoCharge, &vipEndTime)
	mustBinaryRead(buffer, &freshManBonus, &nonoChipList, &dailyRes)
	return UserInfo{
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

// echo测试。服务器原样回复客户端的body
func (c *Connection) Test(parseFunc func(packetBody), data ...interface{}) {
	var id MsgListenerID
	id = c.AddListener(Command_Test, func(body packetBody) {
		c.RemoveListener(Command_Test, id)
		log.Println("Test resp", body.Bytes())
		parseFunc(body)
	})
	err := c.Send(Command_Test, data...)
	if err != nil {
		panic(err)
	}
}
