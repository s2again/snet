package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/fanliao/go-promise"

	"main/config"
	"main/demo/utils"
	"main/snet"
	"main/snet/core"
)

var (
	configFile *config.ServerConfig
	loginAddr  *net.TCPAddr
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
	loginAddr, err = config.GetLoginServer(configFile.IpConfig.HTTP.URL)
	if err != nil {
		panic(err)
	}
	fmt.Println(loginAddr)
}

func main() {
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
		createNewAccount(sid).
			OnSuccess(func(v interface{}) {
				petinfo := v.(snet.PetInfo)
				fmt.Printf("精灵信息：\n%+v\n", petinfo)
				fmt.Printf("个体： %d 性格：%v\n", petinfo.Dv, petinfo.Nature)
			}).
			OnFailure(func(v interface{}) {
				switch v.(type) {
				case error:
					fmt.Println("Error: ", v.(error).Error())
				default:
					fmt.Println("Error: ", v)
				}
			})
	}
}

func createNewAccount(sid string) (task *promise.Promise) {
	task = promise.NewPromise()
	// Login
	_, err := createrolehelper(sid).Get()
	if err != nil {
		task.Reject(errors.New("createrolehelper promise rejected: " + err.Error()))
		return
	}
	loginhelper2(sid).OnSuccess(func(v interface{}) {
		onlineConn := v.(*snet.Connection)
		resp, err := onlineConn.LoginOnline().Get()
		if err != nil {
			task.Reject(err)
		} else {
			task.Resolve(afterlogin(onlineConn, resp.(snet.ResponseForLogin)))
		}
	}).OnFailure(func(v interface{}) {
		task.Reject(v.(error))
	})
	return task
}

// noinspection GoUnusedFunction
func loginhelper(sid string) (prom *promise.Promise) {
	prom = promise.NewPromise()
	defer func() {
		x := recover()
		if x != nil {
			err, ok := x.(error)
			if !ok {
				err = errors.New("promise rejected: " + fmt.Sprint(x))
			}
			prom.Reject(err)
		}
	}()
	uid, sessionID, err := core.ParseSIDString(sid)
	if err != nil {
		panic(err)
	}

	conn, err := snet.Connect(loginAddr)
	if err != nil {
		panic(err)
	}

	conn.SetSession(uid, sessionID)

	conn.ListOnlineServers().OnSuccess(func(v interface{}) {
		// get first online server
		info := v.(snet.CommendSvrInfo)
		fmt.Printf("CommendSvrInfo %+v\n", info)
		server := info.SvrList[0]
		// login online
		addrStr := server.IP + ":" + strconv.Itoa(int(server.Port))
		fmt.Println("Login into Online", addrStr)
		addr, err := net.ResolveTCPAddr("tcp", addrStr)
		if err != nil {
			panic(err)
		}

		onlineConn, err := snet.Connect(addr)
		if err != nil {
			panic(err)
		}

		conn.SetSession(uid, sessionID)
		prom.Resolve(onlineConn)
		conn.Close()
	}).OnFailure(func(v interface{}) {
		conn.Close()
		prom.Reject(v.(error))
	})
	return prom
}

// noinspection GoUnusedFunction
func loginhelper2(sid string) *promise.Promise {
	prom := promise.NewPromise()
	uid, sessionID, err := core.ParseSIDString(sid)
	if err != nil {
		prom.Reject(err)
		return prom
	}

	addr, _ := net.ResolveTCPAddr("tcp4", "182.254.130.223:1222")
	onlineConn, err := snet.Connect(addr)
	if err != nil {
		prom.Reject(err)
		return prom
	}
	onlineConn.SetSession(uid, sessionID)
	prom.Resolve(onlineConn)
	return prom
}

// noinspection GoUnusedFunction
func createrolehelper(sid string) *promise.Promise {
	prom := promise.NewPromise()
	uid, sessionID, err := core.ParseSIDString(sid)
	if err != nil {
		prom.Reject(errors.New("ParseSIDString failed:" + err.Error()))
		return prom
	}

	conn, err := snet.Connect(loginAddr)
	if err != nil {
		prom.Reject(err)
		return prom
	}
	conn.SetSession(uid, sessionID)

	var nickname [16]byte
	copy(nickname[:], "小沙雕")
	_, err = conn.CreateRole(nickname, snet.RoleGreen).Get()
	if err != nil {
		prom.Reject(err)
		return prom
	}
	prom.Resolve(nil)
	return prom
}

func afterlogin(conn *snet.Connection, info snet.ResponseForLogin) snet.PetInfo {
	fmt.Printf("%+v\n", info)
	fmt.Println("登录成功")
	fmt.Printf("userID: %v sessionID: %X\n", conn.UserID, conn.SessionID)
	// Command_SYSTEM_TIME
	// Command_SYSTEM_TIME
	// Command_MAIL_GET_UNREAD
	// Command_NONO_INFO
	utils.AcceptAndCompleteTask(conn, 0x55, 1)               // freshman suit
	task1stPet := utils.AcceptAndCompleteTask(conn, 0x56, 2) // freshman pet
	petTm := task1stPet.CaptureTm
	utils.AcceptAndCompleteTask(conn, 0x57, 1) // freshman pet ball
	utils.AcceptAndCompleteTask(conn, 0x58, 1) // freshman money
	utils.MustResolvePromise(conn.ReleasePet(petTm, 1))
	utils.MustResolvePromise(conn.LeaveMap())
	utils.MustResolvePromise(conn.EnterMap(0, 8, 0x1c8, 0x8f))
	utils.MustResolvePromise(conn.ListMapPlayer())

	petinfo := utils.MustResolvePromise(conn.GetPetInfo(petTm))
	return petinfo.(snet.PetInfo)
}
