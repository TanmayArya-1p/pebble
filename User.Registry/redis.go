package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Connect() *redis.Client {
	var client *redis.Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		Protocol: 2,
	})
	return client
}

func AllOnline(client *redis.Client) []string {
	keys, err := client.Keys(ctx, "*").Result()
	if err != nil {
		fmt.Println("Error getting status:", err)
		return []string{}
	}
	otpt := make([]string, 0)
	for _, key := range keys {
		if strings.HasPrefix(key, "pebble:") {
			otpt = append(otpt, key[7:])
		}
	}
	return otpt
}

func IsOnline(client *redis.Client, uid string) bool {
	_, err := client.Get(ctx, "pebble:"+uid).Result()
	if err == redis.Nil {
		fmt.Println("User ID", uid, "Not in Online DB")
		return false
	} else if err != nil {
		fmt.Errorf("Error getting status:", err)
		return false
	} else {
		return true
	}
}

func SetIsOnline(client *redis.Client, uid string) error {
	td, _ := strconv.Atoi(os.Getenv("USER_ONLINE_FOR"))
	err := client.Set(ctx, "pebble:"+uid, time.Now().UTC().UnixNano(), time.Duration(td)*time.Second).Err()
	if err != nil {
		fmt.Println("Error setting status:", err)
		return err
	}
	return nil
}

func main() {
	godotenv.Load()
	client := Connect()
	defer client.Close()
	uid := "userid123"
	SetIsOnline(client, uid)
	fmt.Println(IsOnline(client, uid))
	fmt.Println(IsOnline(client, "123143"))
	fmt.Println(AllOnline(client))
}
