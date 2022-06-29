package spalkfs

import (
	"context"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func ServeRedisFile(rw http.ResponseWriter, req *http.Request, path string, client *redis.Client) (err error) {
	f, err := client.Get(context.Background(), path).Result()
	if err != nil {
		log.Println(err)
		return
	}
	rw.Write([]byte(f))
	return
}
