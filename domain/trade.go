package domain

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type TradeData struct {
	gorm.Model
	E  string
	E1 int
	S  string
	A  int
	P  string
	Q  string
	F  int
	L  int
	T  int
	M  bool
	M1 bool
}

type TradeApp interface {
	GetLastTradeData() string
	SubTradeWebSocket(connect *websocket.Conn, clientUUID string)
}

type TradeRepo interface {
	GetLastTradeData() string
	SubTradeRedis() <-chan *redis.Message
}
