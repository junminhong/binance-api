package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/junminhong/binance-api/application"
	_ "github.com/junminhong/binance-api/docs"
	"github.com/junminhong/binance-api/domain"
	"github.com/junminhong/binance-api/infra/repo"
	"github.com/junminhong/binance-api/infra/watcher"
	"github.com/junminhong/binance-api/interfaces/http"
	"github.com/junminhong/binance-api/interfaces/http/middleware"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/url"
	"os"
)

var addr = flag.String("addr", "stream.yshyqxx.com", "http service address")

func setUpRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("APP.REDIS_HOST") + ":" + viper.GetString("APP.REDIS_PORT"),
		Password: viper.GetString("APP.REDIS_PASSWORD"), // no password set
		DB:       0,                                     // use default DB
	})
	return client
}

func setUpRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)))
	return router
}

type postgresDB struct {
	db *gorm.DB
}

func setUpDB() *postgresDB {
	// sslmode=disable
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Taipei",
		viper.GetString("APP.DB_HOST"),
		viper.GetString("APP.DB_USERNAME"),
		viper.GetString("APP.DB_PASSWORD"),
		viper.GetString("APP.DB_DATABASE"),
		viper.GetString("APP.DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("Failed to connect DB")
	}
	return &postgresDB{db: db}
}

func (postgresDB *postgresDB) migrationDB() {
	err := postgresDB.db.AutoMigrate(&domain.TradeData{})
	if err != nil {
		log.Println(err.Error())
	}
}

func setUpWebSocket() *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: *addr, Path: "/stream"}
	//?streams=btcusdt@aggTrade
	c, _, err := websocket.DefaultDialer.Dial(u.String()+"?streams=btcusdt@aggTrade", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func setUpDomain(router *gin.Engine, db *gorm.DB, redis *redis.Client, c *websocket.Conn) {
	watcherHandler := watcher.NewHandler(c, db, redis)
	go watcherHandler.WatchTradeData()
	tradeRepo := repo.NewTradeRepo(db, redis)
	tradeApp := application.NewTradeApp(tradeRepo)
	http.NewHandler(router, tradeApp)
}

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
	}
	if viper.GetString("APP.GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// @title           Binance API
// @version         1.0
// @description     好用的東西
// @termsOfService  http://swagger.io/terms/

// @contact.name   junmin.hong
// @contact.url    https://github.com/junminhong
// @contact.email  junminhong1110@gmail.com

// @license.name  MIT
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:8080
// @BasePath  /api/v1
func main() {
	c := setUpWebSocket()
	router := setUpRouter()
	router.Use(middleware.Middleware())
	db := setUpDB()
	go db.migrationDB()
	redis := setUpRedis()
	setUpDomain(router, db.db, redis, c)
	router.Run(viper.GetString("HOST") + ":" + viper.GetString("APP.PORT"))
}
