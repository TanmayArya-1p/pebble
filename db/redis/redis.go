package rd

import (
	"context"
	"fmt"
	"os"
	models "pebble/models"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var clt *redis.Client

func Connect() *redis.Client {
	var client *redis.Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		Protocol: 2,
	})
	clt = client
	return client
}

func AllOnline() []models.UserStatus {
	client := clt
	keys, err := client.Keys(ctx, "*").Result()

	if err != nil {
		fmt.Println("Error getting status:", err)
		return []models.UserStatus{}
	}
	otpt := make([]models.UserStatus, 0)
	for _, key := range keys {
		if strings.HasPrefix(key, "pebble:") {
			curr := models.UserStatus{ID: key[7:], IsOnline: true}
			otpt = append(otpt, curr)
		}
	}
	return otpt
}

func IsOnline(user models.User) models.UserStatus {
	client := clt
	uid := user.ID.Hex()
	_, err := client.Get(ctx, "pebble:"+uid).Result()
	if err == redis.Nil {
		fmt.Println("User ID", uid, "Not Online in DB")
		return models.UserStatus{ID: uid, IsOnline: false}
	} else if err != nil {
		fmt.Println("Error getting status:", err)
		return models.UserStatus{ID: uid, IsOnline: false}
	} else {
		return models.UserStatus{ID: uid, IsOnline: true}
	}
}

func SetIsOnline(user models.User) error {
	client := clt
	uid := user.ID.Hex()
	td, _ := strconv.Atoi(os.Getenv("USER_ONLINE_FOR"))
	err := client.Set(ctx, "pebble:"+uid, time.Now().UTC().UnixNano(), time.Duration(td)*time.Second).Err()
	if err != nil {
		fmt.Println("Error setting status:", err)
		return err
	}
	return nil
}

// func main() {
// 	godotenv.Load()
// 	client := Connect()
// 	defer client.Close()
// 	uid := "userid123"
// 	SetIsOnline(client, uid)
// 	fmt.Println(IsOnline(client, uid))
// 	fmt.Println(IsOnline(client, "123143"))
// 	fmt.Println(AllOnline(client))
// }
