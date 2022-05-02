package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"time"
	//"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	fmt.Println("test")

	accessKey := "xxx"
	accessSecret := "xxxxxxx+"
	secretToken := "xxxxxxxxx"
	stsToken := credentials.NewStaticCredentials(accessKey, accessSecret, secretToken)

	config := aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: stsToken,
	}

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&config))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	filename := "/Users/caoti/Downloads/testFile/wade.jpg"
	f, err := os.Open(filename)
	if err != nil {
		fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	bucket := "ty-us-storage30"
	objectKey := "fa1d2e-69484225-bd126e5fd972e82b/unify/test_wade.jpg"
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
		Body:   f,
	})
	if err != nil {
		fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))

	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
	}
	log.Println("The URL is", urlStr)

}
