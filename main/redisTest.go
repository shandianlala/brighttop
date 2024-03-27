package main

import (
	redis "github.com/go-redis/redis/v7"
	"sdll/brighttop/glog"
	"time"
)

func mainRedis() {

	redisClient := redis.NewClient(
		&redis.Options{
			Network:      "tcp",
			Addr:         "127.0.0.1:1968",
			Password:     "",
			DB:           0,
			PoolSize:     1024,
			MinIdleConns: 100,
			IdleTimeout:  180 * time.Second,
		},
	)

	key1 := "test1"
	//redisClient.Set(key1, "my key...", 30*time.Second)

	result, err := redisClient.Get(key1).Result()
	if err != nil {
		if err != redis.Nil {
			glog.Errorf("获取key出错 %s", err)
		} else {
			glog.Info("key不存在或者已过期")
		}
	} else {
		glog.Infof("value is [%s]", result)
	}
	glog.Info(".....")

	inc, _ := redisClient.Incr("enc_pic_cluster_id").Result()
	glog.Infof("inc=%d", inc)
	intInc, err := redisClient.Get("enc_pic_cluster_id").Int64()
	if err != nil && err != redis.Nil {
		glog.Errorf("获取key出错 %s", err)
	} else {
		glog.Infof("intInc=%d", intInc)
	}

}
