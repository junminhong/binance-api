package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/junminhong/binance-api/domain"
	"github.com/junminhong/binance-api/pkg/responser"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type tradeHandler struct {
	tradeApp domain.TradeApp
}

func NewHandler(router *gin.Engine, tradeApp domain.TradeApp) {
	handler := &tradeHandler{tradeApp: tradeApp}
	router.GET("/api/v1/trade", handler.getLastTradeData)
	router.GET("/api/v1/trade/sub", handler.subTradeWebSocket)
}

// getLastTradeData
// @Summary 取得最新的Trade資料
// @Description
// @Tags trade
// @version 1.0
// @Success 200 {object} aggTradeData "取得成功"
// @Router /trade [get]
func (tradeHandler *tradeHandler) getLastTradeData(c *gin.Context) {
	data := tradeHandler.tradeApp.GetLastTradeData()
	aggTradeData := responser.AggTradeData{}
	json.Unmarshal([]byte(data), &aggTradeData)
	c.JSON(http.StatusOK, responser.Response{
		ResultCode: responser.GetLastTradeDataOk.Code(),
		Message:    responser.GetLastTradeDataOk.Message(),
		Data:       aggTradeData,
		TimeStamp:  time.Now(),
	})
}

func (tradeHandler *tradeHandler) subTradeWebSocket(c *gin.Context) {
	connect, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	clientUUID := uuid.NewString()
	if err != nil {
		log.Println(err.Error())
	}
	tradeHandler.tradeApp.SubTradeWebSocket(connect, clientUUID)
}
