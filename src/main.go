package main

import (
	"log"
	"net/http"

	"github.com/SpalkLtd/video-file-server/src/spalkfs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	if err != nil {
		panic(err.Error())
	}

	svc := s3.New(sess)

	http.Handle("/", spalkfs.FileServer(spalkfs.Dir("./public"), svc, "spalk-video-archive"))

	log.Fatal(http.ListenAndServe(":8663", nil))
}
