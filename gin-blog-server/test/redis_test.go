package test

import (
	"fmt"
	"gin-blog/utils"
)

// var ctx = context.Background()

func ExampleClient() {
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     config.Cfg.Reids.Addr
	// 	Password: config.Cfg.Reids.Password, // no password set
	// 	DB:       config.Cfg.Reids.DB,       // use default DB
	// })

	rdb := utils.InitRedis()
	fmt.Println("redis success: ", rdb)

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }
	// // Output: key value
	// // key2 does not exist
}
