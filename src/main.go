package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/SpalkLtd/video-file-server/src/spalkfs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	gobrake "gopkg.in/airbrake/gobrake.v2"
)

func main() {
	airbrake := gobrake.NewNotifier(1234567, os.Getenv("VFS_ERRBIT_KEY"))
	airbrake.SetHost(os.Getenv("ERRBIT_HOST"))
	defer airbrake.Close()
	defer airbrake.NotifyOnPanic()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION")),
	})
	if err != nil {
		panic(err.Error())
	}

	svc := s3.New(sess)

	http.Handle("/", getFileHandler(spalkfs.FileServer(spalkfs.Dir(os.Getenv("VFS_MEDIA_DIR")), svc, os.Getenv("VFS_S3_BUCKET_FAILOVER")), os.Getenv("VFS_URL_PREFIX")))
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
