package main

import (
	"net/http"
	"os"

	"github.com/SpalkLtd/video-file-server/src/spalkfs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	gobrake "gopkg.in/airbrake/gobrake.v2"
)

func main() {
	airbrake := gobrake.NewNotifier(1234567, "d8b27488dbca7306ad182ff2db2f53d4")
	airbrake.SetHost("https://errbit.spalk.co")
	defer airbrake.Close()
	defer airbrake.NotifyOnPanic()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		panic(err.Error())
	}

	svc := s3.New(sess)

	http.Handle("/", spalkfs.FileServer(spalkfs.Dir("."), svc, "spalk-video-archive"))

	err = http.ListenAndServeTLS("0.0.0.0:443", os.Getenv("CERT_FILE_PATH"), os.Getenv("KEY_FILE_PATH"),nil)
	if err != nil {
		panic(err.Error())
	}
}
