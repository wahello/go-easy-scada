/**
 * @Author: cmpeax Tang
 * @Date: 2021/2/25 19:43
 */

package wsckio

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
)

type SocketIOService struct {
	addr    string
	sockets sync.Map
	handler map[string]func(socketio.Conn, interface{}) string
}

func NewSocketIOService(addr string) *SocketIOService {
	return &SocketIOService{
		addr: addr,
	}
}

func (t *SocketIOService) Mount(handler map[string]func(socketio.Conn, interface{}) string) {
	t.handler = handler
}

func (t *SocketIOService) Init(cnncb func(s socketio.Conn)) error {
	router := gin.New()

	server, err := socketio.NewServer(&engineio.Options{})
	if err != nil {
		return err
	}

	server.OnConnect("/", func(s socketio.Conn) error {

		t.sockets.Store(s.RemoteAddr(), s)
		cnncb(s)
		return nil
	})

	for k, v := range t.handler {
		server.OnEvent("/", k, v)
	}

	// server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
	// 	s.SetContext(msg)
	// 	return "recv " + msg
	// })

	// server.OnEvent("/", "bye", func(s socketio.Conn) string {
	// 	last := s.Context().(string)
	// 	s.Emit("bye", last)
	// 	s.Close()
	// 	return last
	// })

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	// defer server.Close()

	router.Use(ginMiddleware("*"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.Run(t.addr)
	return nil
}

// func NewSocketIO() {
// 	router := gin.New()

// 	server, err := socketio.NewServer(&engineio.Options{})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	server.OnConnect("/", func(s socketio.Conn) error {

// 		fmt.Println("connected:", s.ID())
// 		return nil
// 	})

// 	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
// 		fmt.Println("notice:", msg)
// 		s.Emit("reply", "have "+msg)
// 	})

// 	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
// 		s.SetContext(msg)
// 		return "recv " + msg
// 	})

// 	server.OnEvent("/", "bye", func(s socketio.Conn) string {
// 		last := s.Context().(string)
// 		s.Emit("bye", last)
// 		s.Close()
// 		return last
// 	})

// 	server.OnError("/", func(s socketio.Conn, e error) {
// 		fmt.Println("meet error:", e)
// 	})

// 	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// 		fmt.Println("closed", reason)
// 	})

// 	go server.Serve()
// 	// defer server.Close()

// 	router.Use(ginMiddleware("*"))
// 	router.GET("/socket.io/*any", gin.WrapH(server))
// 	router.POST("/socket.io/*any", gin.WrapH(server))
// 	router.Run(":9994")
// 	// http.Handle("/socket.io/", server)
// 	// http.ListenAndServe(":9994", nil)
// }

func ginMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}
