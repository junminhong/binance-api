package repo

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/junminhong/binance-api/domain"
	"gorm.io/gorm"
)

type tradeRepo struct {
	redis *redis.Client
	db    *gorm.DB
}

func (tradeRepo *tradeRepo) GetLastTradeData() string {
	return tradeRepo.redis.Get(context.Background(), "streams=btcusdt@aggTrade").Val()
}

func (tradeRepo *tradeRepo) SubTradeRedis() <-chan *redis.Message {
	pubsub := tradeRepo.redis.Subscribe(context.Background(), "trade_sub")
	ch := pubsub.Channel()
	return ch
}

func NewTradeRepo(db *gorm.DB, redis *redis.Client) domain.TradeRepo {
	return &tradeRepo{redis: redis, db: db}
}
