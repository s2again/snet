package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/fanliao/go-promise"

	"main/connection"
	"main/connection/core"
)

func main_reghelper() {
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
				fmt.Printf("个体： %d\n", v)
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
	select {}
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
		onlineConn := v.(*connection.Connection)
		resp, err := onlineConn.LoginOnline().Get()
		if err != nil {
			task.Reject(err)
		} else {
			task.Resolve(afterlogin(onlineConn, resp.(connection.ResponseForLogin)))
		}
	}).OnFailure(func(v interface{}) {
		task.Reject(v.(error))
	})
	return task
}

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

	conn, err := connection.Connect(loginAddr)
	if err != nil {
		panic(err)
	}

	conn.SetSession(uid, sessionID)

	conn.ListOnlineServers().OnSuccess(func(v interface{}) {
		// get first online server
		info := v.(connection.CommendSvrInfo)
		fmt.Printf("CommendSvrInfo %+v\n", info)
		server := info.SvrList[0]
		// login online
		addrStr := server.IP + ":" + strconv.Itoa(int(server.Port))
		fmt.Println("Login into Online", addrStr)
		addr, err := net.ResolveTCPAddr("tcp", addrStr)
		if err != nil {
			panic(err)
		}

		onlineConn, err := connection.Connect(addr)
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

func loginhelper2(sid string) *promise.Promise {
	prom := promise.NewPromise()
	uid, sessionID, err := core.ParseSIDString(sid)
	if err != nil {
		prom.Reject(err)
		return prom
	}

	addr, _ := net.ResolveTCPAddr("tcp4", "182.254.130.223:1222")
	onlineConn, err := connection.Connect(addr)
	if err != nil {
		prom.Reject(err)
		return prom
	}
	onlineConn.SetSession(uid, sessionID)
	prom.Resolve(onlineConn)
	return prom
}

func createrolehelper(sid string) *promise.Promise {
	prom := promise.NewPromise()
	uid, sessionID, err := core.ParseSIDString(sid)
	if err != nil {
		prom.Reject(errors.New("ParseSIDString failed:" + err.Error()))
		return prom
	}

	conn, err := connection.Connect(loginAddr)
	if err != nil {
		prom.Reject(err)
		return prom
	}
	conn.SetSession(uid, sessionID)

	var nickname [16]byte
	copy(nickname[:], "小沙雕")
	_, err = conn.CreateRole(nickname, connection.RoleGreen).Get()
	if err != nil {
		prom.Reject(err)
		return prom
	}
	prom.Resolve(nil)
	return prom
}

func afterlogin(conn *connection.Connection, info connection.ResponseForLogin) uint32 {
	fmt.Printf("%+v\n", info)
	fmt.Println("登录成功")
	fmt.Printf("userID: %v sessionID: %X\n", conn.UserID, conn.SessionID)
	// Command_SYSTEM_TIME
	// Command_SYSTEM_TIME
	// Command_MAIL_GET_UNREAD
	// Command_NONO_INFO

	acceptAndCompleteTask(conn, 0x55, 1)               // freshman suit
	task1stPet := acceptAndCompleteTask(conn, 0x56, 2) // freshman pet
	petTm := task1stPet.CaptureTm
	acceptAndCompleteTask(conn, 0x57, 1) // freshman pet ball
	acceptAndCompleteTask(conn, 0x58, 1) // freshman money
	mustResolvePromise(conn.ReleasePet(petTm, 1))
	mustResolvePromise(conn.LeaveMap())
	mustResolvePromise(conn.EnterMap(0, 8, 0x1c8, 0x8f))
	mustResolvePromise(conn.ListMapPlayer())

	petinfo := mustResolvePromise(conn.GetPetInfo(petTm))
	return petinfo.(connection.PetInfo).Dv

}
