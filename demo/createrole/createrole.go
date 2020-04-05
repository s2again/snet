package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fanliao/go-promise"

	"main/config"
	"main/demo/utils"
	"main/snet"
)

var (
	configFile *config.ServerConfig
	guideAddr  *net.TCPAddr
)

func init() {
	var err error
	f, err := os.OpenFile("seer.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	log.SetOutput(f)
	configFile, err = config.GetServerConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(configFile)
	guideAddr, err = configFile.GetGuideServerByHTTP()
	if err != nil {
		panic(err)
	}
	fmt.Println(guideAddr)
}

func main() {
	fmt.Print("输入新手宠物任务参数:")
	var noviceParam1 uint32
	n, err := fmt.Scanf("%d\n", &noviceParam1)
	if n < 1 {
		fmt.Println("too few input")
		os.Exit(-1)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	for {
		var sid string
		fmt.Print("输入SID:")
		n, err := fmt.Scanf("%40s\n", &sid)
		if n < 1 {
			fmt.Println("too few input")
			continue
		}
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		v, err := createFreshmenRole(sid, noviceParam1).Get()
		if err != nil {
			log.Println(err)
			fmt.Println("Error: ", err)
			continue
		}
		petinfo := v.(snet.PetInfo)
		fmt.Printf("精灵信息：\n%+v\n", petinfo)
		fmt.Printf("个体： %d 性格：%v\n", petinfo.Dv, petinfo.Nature)
	}
}

func createFreshmenRole(sid string, noviceParam1 uint32) (task *promise.Promise) {
	task = promise.NewPromise()
	go func() {
		// connect guide server
		v, err := utils.ConnectGuideServer(guideAddr, sid).Get()
		if err != nil {
			task.Reject(errors.New("ConnectGuideServer promise rejected: " + err.Error()))
			return
		}
		// create
		conn := v.(*snet.GuideServerConnection)
		v, err = createRole(conn).Get()
		if err != nil {
			task.Reject(errors.New("createRole promise rejected: " + err.Error()))
			return
		}

		// get online server list
		list, err := utils.GetOnlineServerList(conn)
		if err != nil || len(list) < 1 {
			task.Reject(errors.New("getServerList failed: " + err.Error()))
			return
		}
		// connect first online server, and close connection to sub server
		v, err = utils.Guide2Online(conn, list[0]).Get()
		if err != nil {
			task.Reject(errors.New("connectOnline promise rejected: " + err.Error()))
			return
		}
		onlineConn := v.(*snet.Connection)
		// login online server
		resp, err := onlineConn.LoginOnline().Get()
		task.Resolve(finishNoviceAfterLogin(onlineConn, resp.(snet.ResponseForLogin), noviceParam1))
		onlineConn.Close()
	}()
	return task
}

func finishNoviceAfterLogin(conn *snet.OnlineServerConnection, info snet.ResponseForLogin, noviceParam1 uint32) (pet snet.PetInfo, err error) {
	defer func() {
		x := recover()
		if x != nil {
			err = errors.New("Error: " + fmt.Sprint(x))
			log.Println(x)
			fmt.Println("发生错误：", x)
		}
	}()
	fmt.Printf("%+v\n", info)
	fmt.Println("登录成功")
	fmt.Printf("userID: %v sessionID: %X\n", conn.UserID, conn.SessionID)
	// Command_SYSTEM_TIME
	// Command_SYSTEM_TIME
	// Command_MAIL_GET_UNREAD
	// Command_NONO_INFO
	utils.AcceptAndCompleteTask(conn, 0x55, 1)                          // freshman suit
	task1stPet := utils.AcceptAndCompleteTask(conn, 0x56, noviceParam1) // freshman pet
	petTm := task1stPet.CaptureTm
	utils.AcceptAndCompleteTask(conn, 0x57, 1) // freshman pet ball
	utils.AcceptAndCompleteTask(conn, 0x58, 1) // freshman money
	utils.MustResolvePromise(conn.ReleasePet(petTm, 1))
	utils.MustResolvePromise(conn.LeaveMap())
	utils.MustResolvePromise(conn.EnterMap(0, 8, 0x1c8, 0x8f))
	utils.MustResolvePromise(conn.ListMapPlayer())

	petinfo := utils.MustResolvePromise(conn.GetPetInfo(petTm))
	return petinfo.(snet.PetInfo), nil
}

// noinspection GoUnusedFunction
func createRole(conn *snet.GuideServerConnection) *promise.Promise {
	prom := promise.NewPromise()
	go func() {
		var nickname [16]byte
		copy(nickname[:], "小沙雕")
		_, err := conn.CreateRole(nickname, snet.RoleGreen).Get()
		if err != nil {
			prom.Reject(err)
			return
		}
		prom.Resolve(nil)
	}()
	return prom
}
