package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/SpalkLtd/video-file-server/spalkfs"
	"github.com/airbrake/gobrake/v5"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis/v8"
)

func main() {
	airbrake := gobrake.NewNotifierWithOptions(&gobrake.NotifierOptions{
		ProjectId:  1234567,
		ProjectKey: os.Getenv("VFS_ERRBIT_KEY"),
		Host:       os.Getenv("ERRBIT_HOST"),
	})
	defer airbrake.NotifyOnPanic()
	defer airbrake.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION")),
	})
	if err != nil {
		panic(err.Error())
	}

	svc := s3.New(sess)

	mediaDir := os.Getenv("VFS_MEDIA_DIR")
	if mediaDir == "" || mediaDir == "." || mediaDir == ".." {
		mediaDir = "public"
	}
	fi, err := os.Stat(mediaDir)
	if err != nil {
		log.Println("Error checking media dir")
		log.Println(err.Error())
		return
	}
	if !fi.IsDir() {
		log.Println("VFS_MEDIA_DIR must be set to a directory")
		return
	}

	redisTTL, _ := strconv.Atoi(os.Getenv("MS_REDIS_TTL"))
	if redisTTL == 0 {
		redisTTL = 20
	}

	rClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // Default
	})

	http.Handle("/", getFileHandler(spalkfs.FileServer(spalkfs.Dir(mediaDir), rClient, svc, os.Getenv("VFS_S3_BUCKET_FAILOVER")), os.Getenv("VFS_URL_PREFIX")))
	certpath, keypath := os.Getenv("VFS_CERT_FILE_PATH"), os.Getenv("VFS_KEY_FILE_PATH")
	if certpath == "" || keypath == "" {
		fmt.Println("insufficient signing info found. Defaulting to http on localhost")
		err = http.ListenAndServe(os.Getenv("VFS_HTTP_BIND_ADDRESS"), nil)
	} else {
		err = http.ListenAndServeTLS(os.Getenv("VFS_HTTPS_BIND_ADDRESS"), os.Getenv("VFS_CERT_FILE_PATH"), os.Getenv("VFS_KEY_FILE_PATH"), nil)
	}
	if err != nil {
		panic(err.Error())
	}
}

//make the standard library file server bahave with gocraft/web
func getFileHandler(handler spalkfs.Handler, prefix string) http.Handler {
	if prefix != "" {
		fmt.Printf("stripping prefix: %v\n", prefix)
		return http.StripPrefix(prefix, handler)
	}
	return handler
}
