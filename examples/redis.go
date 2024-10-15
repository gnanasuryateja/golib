package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/gnanasuryateja/golib/datastore/cache/redis"
// )

// type s struct {
// 	Id   string
// 	Data any
// }

// func main() {
// 	ctx := context.Background()
// 	redisStoreConfig := redis.RedisStoreConfig{
// 		Addr:     "",
// 		Port:     "",
// 		Username: "",
// 		Password: "",
// 	}
// 	cache, err := redis.NewRedisStoreClient(ctx, redisStoreConfig)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// testing cache.HealthCheck
// 	err = cache.HealthCheck(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	} else {
// 		fmt.Println("health check successful...:)")
// 	}
// 	s1 := s{
// 		Id:   "s1",
// 		Data: "Sample struct 1",
// 	}
// 	// testing cache.AddData
// 	resAddData, err := cache.AddData(ctx, s1.Id, s1)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("resAddData:", resAddData)
// 	// testing cache.GetData
// 	resGetData, err := cache.GetData(ctx, s1.Id)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("resGetData:", string(resGetData.([]byte)))
// 	// testing cache.GetKeys
// 	keys, err := cache.GetKeys(ctx, "")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("keys:", keys)
// 	// testing cache.DeleteData
// 	resDeleteData, err := cache.DeleteData(ctx, "s1")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("resDeleteData:", resDeleteData)
// }
