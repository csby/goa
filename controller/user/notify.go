package user

import (
	"encoding/json"
	"github.com/csby/goa/controller"
	"github.com/csby/goa/data/socket"
	"github.com/csby/gwsf/gtype"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

func NewNotify(log gtype.Log, param *controller.Parameter) *Notify {
	instance := &Notify{}
	instance.SetLog(log)
	instance.SetParameter(param)

	instance.wsGrader = websocket.Upgrader{CheckOrigin: instance.checkOrigin}

	if instance.WChs != nil {
		instance.WChs.SetListener(nil, instance.onChannelRemoved)
		//instance.WChs.AddReader(instance.onChannelRead)
	}

	return instance
}

type Notify struct {
	base

	wsGrader websocket.Upgrader
}

func (s *Notify) Socket(ctx gtype.Context, ps gtype.Params) {
	websocketConn, err := s.wsGrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		s.LogError("login socket connect fail:", err)
		ctx.Error(gtype.ErrInternal, err)
		return
	}
	defer websocketConn.Close()

	token := s.GetToken(ctx.Token())
	if token != nil {
		s.Tdb.Permanent(token.ID, true)

		s.WChs.Write(&gtype.SocketMessage{
			ID: socket.WSUserLogin,
			Data: &gtype.OnlineUser{
				UserAccount: token.UserAccount,
				UserName:    token.UserName,
				LoginIP:     token.LoginIP,
				LoginTime:   gtype.DateTime(time.Now()),
			},
		}, token)
	}
	channel := s.WChs.NewChannel(token)
	defer s.WChs.Remove(channel)

	waitGroup := &sync.WaitGroup{}
	stopWrite := make(chan bool, 2)
	stopRead := make(chan bool, 2)

	// write message
	waitGroup.Add(1)
	go func(wg *sync.WaitGroup, conn *websocket.Conn, ch gtype.SocketChannel) {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				s.LogError("login socket send message error:", err)
			}
			stopRead <- true
		}()

		for {
			select {
			case <-stopWrite:
				return
			case msg, ok := <-ch.Read():
				if !ok {
					return
				}

				conn.WriteJSON(msg)
			}
		}
	}(waitGroup, websocketConn, channel)

	// read message
	waitGroup.Add(1)
	go func(wg *sync.WaitGroup, conn *websocket.Conn, ch gtype.SocketChannel) {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				s.LogError("login socket send message error:", err)
			}
			stopWrite <- true
		}()

		for {
			select {
			case <-stopRead:
				return
			default:
				msgType, msgContent, err := conn.ReadMessage()
				if err != nil {
					return
				}
				if msgType == websocket.CloseMessage {
					return
				}

				if msgType == websocket.TextMessage || msgType == websocket.BinaryMessage {
					msg := &gtype.SocketMessage{}
					err := json.Unmarshal(msgContent, msg)
					if err == nil {
						s.WChs.Read(msg, ch)
					}
				}
			}
		}
	}(waitGroup, websocketConn, channel)

	waitGroup.Wait()
}

func (s *Notify) SocketDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, userCatalogLogin)
	function := catalog.AddFunction(method, uri, "消息推送")
	function.SetNote("订阅并接收系统推送的消息，该接口保持阻塞至连接关闭")
	function.SetInputExample(&gtype.SocketMessage{ID: 1})
	function.SetOutputExample(&gtype.SocketMessage{ID: 1})
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrTokenInvalid)
}

func (s *Notify) checkOrigin(r *http.Request) bool {
	if r != nil {
	}
	return true
}

func (s *Notify) onChannelRemoved(channel gtype.SocketChannel) {
	if channel == nil {
		return
	}

	token := channel.Token()
	if token == nil {
		return
	}

	if token.Usage > 0 {
		return
	}

	if s.Tdb != nil {
		s.Tdb.Permanent(token.ID, false)
	}
}

func (s *Notify) onChannelRead(message *gtype.SocketMessage, channel gtype.SocketChannel) {
	channel.Container().Write(&gtype.SocketMessage{
		ID:   message.ID,
		Data: message.Data,
	}, channel.Token())
}
