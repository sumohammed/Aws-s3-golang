package fileuploader

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
)

func Uploader(id string, key string, token string, data []byte, name string) {
	size := int64(len(data))
	log.Println(name)
	buffer := make([]byte, size) // read file content to buffer

	copy(buffer, data)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	log.Println(fileType)

	creds := credentials.NewStaticCredentials(id, key, token)
	_, err := creds.Get()
	if err != nil {
		log.Printf("bad credentials: %s", err)
	}
	cfg := aws.NewConfig().WithRegion("eu-west-2").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	path := "/farms/" + name
	params := &s3.PutObjectInput{
		Bucket:        aws.String("farmcap"),
		Key:           aws.String(path),
		ACL:           aws.String("public-read"),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		log.Printf("bad response: %s", err)
	}
	log.Printf("response %s", awsutil.StringValue(resp))
}
