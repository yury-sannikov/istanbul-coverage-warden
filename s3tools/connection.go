package s3tools

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// S3ConnectionOptions hold connection options to connect to S3s
type S3ConnectionOptions struct {
	Region *string
}

// CreateS3Session return new s3.S3 object
func CreateS3Session(options *S3ConnectionOptions) (*session.Session, error) {
	config := &aws.Config{Region: options.Region}
	return session.New(config), nil
}
