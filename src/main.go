package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	// provider := ec2rolecreds.EC2RoleProvider{}
	// creds, err := provider.Retrieve()
	// if err != nil {
	// 	panic(err.Error())
	// }

	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("ap-southeast-2a"),
	// 	Credentials: credentials.NewStaticCredentialsFromCreds(creds),
	// })
	// if err != nil {
	// 	panic(err.Error())
	// }

	// svc := s3.New(sess)

	// http.Handle("/", FileServer(Dir("public"), svc))
	tmp := s3.S3{}
	http.Handle("/", FileServer(Dir("public"), &tmp))

	log.Fatal(http.ListenAndServe(":8090", nil))
}
