package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func serveS3File(rw http.ResponseWriter, req *http.Request, name string, s3svc *s3.S3) {

	params := s3.GetObjectInput{
		Bucket: aws.String("spalk-video-archive"),
		Key:    aws.String(name),
	}

	resp, err := s3svc.GetObject(&params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	h := rw.Header()
	h.Add("content-type", *resp.ContentType)
	_, err = io.Copy(rw, resp.Body)
}
