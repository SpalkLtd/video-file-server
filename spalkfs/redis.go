package spalkfs

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

func ServeRedisObject(rw http.ResponseWriter, req *http.Request, name string, redisClient *redis.Client, prefix string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	vars := strings.Split(name)
	var rKey string
	if len(vars) == 4 {
		rKey = fmt.Sprintf("%s:%s:%s:%s:%s", os.Getenv("REDIS_PREFIX"), vars[0], vars[1], vars[2], vars[3])
	} else if len {
		rKey = fmt.Sprintf("%s:%s:%s:%s", os.Getenv("REDIS_PREFIX"), vars[0], vars[1], vars[2])
	}

	rkey = prefix + strings.Replace(name, "/", ":", -1)

	pl, err := rClient.Get(rKey).Result()
	if err == redis.Nil {
		// Not in Redis, do S3 download

		// Update the cache with the playlist
		err = rClient.Set(rKey, pl, 0).Err()

		if err != nil {
			glog.Errorf("Error putting playlist into cache: %s", err)
		}
	}

	buff := new(bytes.Buffer)
	// GetObject returns a buffer, read it into a string
	buff.ReadFrom(result.Body)
	pl := buff.String()
	result.Body.Close()

}
