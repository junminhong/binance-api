package watcher

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/junminhong/binance-api/domain"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"sync"
)

type data struct {
	E  string `json:"e"`
	E1 int    `json:"E"`
	S  string `json:"s"`
	A  int    `json:"a"`
	P  string `json:"p"`
	Q  string `json:"q"`
	F  int    `json:"f"`
	L  int    `json:"l"`
	T  int    `json:"T"`
	M  bool   `json:"m"`
	M1 bool   `json:"M"`
}

type aggTradeData struct {
	Stream string
	Data   data
}

type watcherHandler struct {
	c     *websocket.Conn
	redis *redis.Client
	db    *gorm.DB
}

func NewHandler(c *websocket.Conn, db *gorm.DB, redis *redis.Client) *watcherHandler {
	return &watcherHandler{c: c, db: db, redis: redis}
}

func (watcherHandler *watcherHandler) WatchTradeData() {
	var tradeDataList []domain.TradeData
	for {
		_, message, err := watcherHandler.c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		aggTradeData := aggTradeData{}
		json.Unmarshal(message, &aggTradeData)
		tradeData := domain.TradeData{
			E:  aggTradeData.Data.E,
			E1: aggTradeData.Data.E1,
			S:  aggTradeData.Data.S,
			A:  aggTradeData.Data.A,
			P:  aggTradeData.Data.P,
			Q:  aggTradeData.Data.Q,
			F:  aggTradeData.Data.F,
			L:  aggTradeData.Data.L,
			T:  aggTradeData.Data.T,
			M:  aggTradeData.Data.M,
			M1: aggTradeData.Data.M1,
		}
		// log.Println(time.Unix(int64(aggTradeData.Data.E1), 0))
		wg := sync.WaitGroup{}
		wg.Add(1)
		watcherHandler.updateTradeData(message, &wg)
		wg.Wait()
		if viper.GetBool("APP.SAVE_PG") {
			tradeDataList = append(tradeDataList, tradeData)
			if len(tradeDataList) == viper.GetInt("APP.SAVE_PG_MAX") {
				wg := sync.WaitGroup{}
				wg.Add(1)
				go watcherHandler.saveTradeData(&tradeDataList, &wg)
				wg.Wait()
			}
		}
	}
}

func (watcherHandler *watcherHandler) updateTradeData(data []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	err := watcherHandler.redis.Set(context.Background(), "streams=btcusdt@aggTrade", data, 0).Err()
	if err != nil {
		log.Println(err.Error())
	}
	err = watcherHandler.redis.Publish(context.Background(), "trade_sub", data).Err()
	if err != nil {
		log.Println(err.Error())
	}
}

func (watcherHandler *watcherHandler) saveTradeData(data *[]domain.TradeData, wg *sync.WaitGroup) {
	defer wg.Done()
	watcherHandler.db.CreateInBatches(&data, viper.GetInt("APP.SAVE_PG_MAX"))
	*data = nil
}
