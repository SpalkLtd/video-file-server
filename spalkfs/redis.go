package spalkfs

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

func ServeRedisFile(rw http.ResponseWriter, req *http.Request, path string, client *redis.Client) (err error) {
	// Prefix was stripped, need to add back in
	prefix := os.Getenv("LH_FS_URL_PREFIX")
	f, err := client.Get(prefix + path).Result()
	if err != nil {
		log.Println(err)
		return
	}
	rw.Write([]byte(f))
	return
}
