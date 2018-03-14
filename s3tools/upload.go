package s3tools

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadFile uploads file to the specific S3 bucket
func UploadFile(sess *session.Session, bucketName string, uploadPath string, filename string) error {

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uploadPath),
		Body:   f,
	})
	if err != nil {
		return err
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)

	return nil
}
