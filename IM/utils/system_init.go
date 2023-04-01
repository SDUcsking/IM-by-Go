package utils

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config mysql", viper.Get("mysql"))
}

func InitRedis() {
	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

}

func InitMySql() {
	newLogger := logger.New(
		//自定义日志模板,打印sql语句
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢sql阈值
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)
	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	fmt.Println("mysql inited")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}

const (
	PublishKey = "websocket"
)

// 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	fmt.Println("Publish...", msg)
	err = Red.Publish(ctx, channel, msg).Err()
	return err
}

//订阅消息到redis

func Subscribe(ctx context.Context, channel string) (string, error) {

	sub := Red.Subscribe(ctx, channel)
	fmt.Println("Subscribe", ctx)

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println("就是这里,", err)
	}
	fmt.Printf("这一步")
	fmt.Printf(msg.Payload)
	//fmt.Println("subscribe", msg.Payload)
	fmt.Printf("这一步")

	return msg.Payload, err
}
