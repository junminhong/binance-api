package application

import (
	"github.com/gorilla/websocket"
	"github.com/junminhong/binance-api/domain"
	"log"
	"sync"
)

var (
	mux    sync.Mutex
	client = make(map[string]*websocket.Conn)
)

type tradeApp struct {
	tradeRepo domain.TradeRepo
}

func (tradeApp *tradeApp) SubTradeWebSocket(connect *websocket.Conn, clientUUID string) {
	addClient(clientUUID, connect)
	ch := tradeApp.tradeRepo.SubTradeRedis()
	defer connect.Close()
	for message := range ch {
		if err := connect.WriteMessage(websocket.TextMessage, []byte("最新一筆資料: "+message.Payload)); err != nil {
			log.Println(err.Error())
		}
	}
}

func (tradeApp *tradeApp) GetLastTradeData() string {
	return tradeApp.tradeRepo.GetLastTradeData()
}

func NewTradeApp(tradeRepo domain.TradeRepo) domain.TradeApp {
	return &tradeApp{tradeRepo: tradeRepo}
}

func addClient(clientUUID string, conn *websocket.Conn) {
	mux.Lock()
	client[clientUUID] = conn
	mux.Unlock()
}
