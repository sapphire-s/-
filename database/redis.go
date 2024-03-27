package database

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func GetRedis() *redis.Client {
	addr := viper.Get("redis.addr").(string)
	passwd := viper.Get("redis.passwd").(string)
	db := int(viper.Get("redis.db").(float64))
	fmt.Println("redis_parse:", "addr:"+addr, "passwd:"+passwd, "db:", db)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
}
