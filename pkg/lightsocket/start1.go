package lightsocket

//
//import (
//	"github.com/gobwas/ws"
//	"github.com/gobwas/ws/wsutil"
//	"github.com/gofiber/fiber/v2"
//	"golang.org/x/sys/unix"
//	"log"
//)
//
//var epoller *epoll
//
//func WSHandler(ctx *fiber.Ctx) {
//	conn := ctx.Context().Conn()
//	_, err := ws.Upgrade(conn)
//	if err != nil {
//		log.Printf("Failed to upgrade connection %v", err)
//		conn.Close()
//	}
//	if err = epoller.Add(conn); err != nil {
//		log.Printf("Failed to add connection %v", err)
//		conn.Close()
//	}
//}
//
//func Init() {
//	// Increase resources limitations
//	var rLimit unix.Rlimit
//	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &rLimit); err != nil {
//		panic(err)
//	}
//	rLimit.Cur = rLimit.Max
//	if err := unix.Setrlimit(unix.RLIMIT_NOFILE, &rLimit); err != nil {
//		panic(err)
//	}
//
//	// Start epoll
//	var err error
//	epoller, err = MkEpoll()
//	if err != nil {
//		panic(err)
//	}
//
//	go Start()
//}
//
//func Start() {
//	for {
//		connections, err := epoller.Wait()
//		if err != nil {
//			log.Printf("Failed to epoll wait %v", err)
//			continue
//		}
//		for _, conn := range connections {
//			if conn == nil {
//				break
//			}
//			var message []byte
//			var opCode ws.OpCode
//			if message, opCode, err = wsutil.ReadClientData(conn); err != nil {
//				log.Println(err)
//				if err = epoller.Remove(conn); err != nil {
//					log.Printf("Failed to remove %v", err)
//				}
//				conn.Close()
//			} else {
//				log.Println(message, opCode)
//			}
//		}
//	}
//}
