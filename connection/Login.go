package connection

import (
	"crypto/md5"
	"fmt"
	"log"

	"github.com/fanliao/go-promise"

	"main/connection/core"
)

func (c *Connection) LoginByEmail(email string, password string) (prom *promise.Promise) {
	const channel uint32 = 30
	const gameType uint32 = 2
	t := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	t = fmt.Sprintf("%x", md5.Sum([]byte(t)))
	var pwd [32]byte
	copy(pwd[:], []byte(t)[:32])
	log.Println(t)
	var emailBytes [64]byte
	copy(emailBytes[:], []byte(email)[:64])

	prom = promise.NewPromise()
	c.SetSession(0, [16]byte{})
	c.SendInPromise(Command_MAIN_LOGIN_IN, emailBytes, pwd, channel, gameType, uint32(0)).
		OnSuccess(func(v interface{}) {
			fmt.Printf("LoginResponse %X\n", v.(core.PacketBody).Bytes())
			prom.Resolve(v)
		}).
		OnFailure(func(v interface{}) {
			prom.Reject(v.(error))
		})
	return prom
}

func (c *Connection) Login(password string) (prom *promise.Promise) {
	const channel uint32 = 30
	const gameType uint32 = 2
	t := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	t = fmt.Sprintf("%x", md5.Sum([]byte(t)))
	var pwd [32]byte
	copy(pwd[:], []byte(t)[:32])
	log.Println(t)
	prom = promise.NewPromise()

	/*c.SendInPromise(103, pwd, channel, gameType, uint32(0), [16]byte{}, [6]byte{}, [64]byte{byte('0'), 0}).
	OnSuccess(func(v interface{}) {
		fmt.Printf("LoginResponse %X\n", v.(core.PacketBody).Bytes())
		prom.Resolve(v)
	}).
	OnFailure(func(v interface{}) {
		prom.Reject(v.(error))
	})
	*/

	c.SendInPromise(Command_MAIN_LOGIN_IN, pwd, channel, gameType, uint32(0)).
		OnSuccess(func(v interface{}) {
			fmt.Printf("LoginResponse %X\n", v.(core.PacketBody).Bytes())
			prom.Resolve(v)
		}).
		OnFailure(func(v interface{}) {
			prom.Reject(v.(error))
		})
	return prom
}
