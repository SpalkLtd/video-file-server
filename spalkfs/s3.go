package spalkfs

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func ServeS3File(rw http.ResponseWriter, req *http.Request, name string, s3svc *s3.S3, bucket string) (err error) {
	bucketParts := strings.Split(bucket, "/")
	bucketName := bucketParts[0]
	bucketPath := strings.Join(bucketParts[1:], "/")
	log.Println(bucketPath + name)
	params := s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bucketPath + name),
	}

	resp, err := s3svc.GetObject(&params)

	if err != nil {
		log.Println(err)
		return
	}

	h := rw.Header()
	if !strings.HasSuffix(name, "m3u8") {
		cacheUntil := time.Now().AddDate(0, 0, 30).Format(http.TimeFormat)
		h.Set("Expires", cacheUntil)
		h.Set("Cache-Control", "max-age=2592000")
	}
	h.Add("content-type", *resp.ContentType)
	_, err = io.Copy(rw, resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func isS3NotFound(err error) bool {
	return err.Error()[:10] == "NoSuchKey:" || err.Error()[:22] == "NoCredentialProviders:"
}
