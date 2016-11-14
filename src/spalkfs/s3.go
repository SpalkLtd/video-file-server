package spalkfs

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ServeS3File(rw http.ResponseWriter, req *http.Request, name string, s3svc *s3.S3, bucket string) {

	params := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(name),
	}

	resp, err := s3svc.GetObject(&params)

	if err != nil {
		fmt.Println(err.Error())
		if isS3NotFound(err) {
			http.Error(rw, "Not Found", http.StatusNotFound)
		} else {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	h := rw.Header()
	if !strings.HasSuffix(name, "m3u8") {
		h.Set("Cache-Control", "max-age=2592000")
	}
	h.Add("content-type", *resp.ContentType)
	_, err = io.Copy(rw, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func isS3NotFound(err error) bool {
	return err.Error()[:10] == "NoSuchKey:"
}
