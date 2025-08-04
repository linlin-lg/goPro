package Connect

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connect struct {
	Conn     *websocket.Conn
	DataChan chan interface{}
}

var connect Connect

func InitConnect(c *websocket.Conn, ch chan interface{}) Connect {
	connect = Connect{c, ch}

	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			fmt.Println("Conn ReadMessage err = ", err)
		} else {
			println("ReadMessage = ", data)
		}
	}

	return connect
}

func (c *Connect) SendMessage(code int, data []byte) {
	c.Conn.WriteMessage(code, data)
}
