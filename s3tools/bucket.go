package s3tools

import (
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const _NotFoundError = "NotFound"

// CheckOrCreateBucket check if specified bucket exists, if not, creates new one
func CheckOrCreateBucket(session *session.Session, bucketName string) error {

	service := s3.New(session)

	_, err := service.HeadBucket(&s3.HeadBucketInput{Bucket: &bucketName})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == _NotFoundError {
				fmt.Printf("Bucket '%s' does not exists\n", bucketName)
				return createBucket(service, bucketName)
			}
		}
		log.Printf("Error %s while checking bucket %s: %+v\n", reflect.TypeOf(err), bucketName, err)
		return err
	}

	return nil
}

func createBucket(service *s3.S3, bucketName string) error {
	result, err := service.CreateBucket(&s3.CreateBucketInput{
		Bucket: &bucketName,
	})

	if err != nil {
		log.Printf("Failed to create bucket %s, %+v\n", bucketName, err)
		return err
	}

	if err = service.WaitUntilBucketExists(&s3.HeadBucketInput{Bucket: &bucketName}); err != nil {
		log.Printf("Failed to wait for bucket to exist %s, %s\n", bucketName, err)
		return err
	}

	log.Printf("Successfully created bucket %s\n", result)

	policyErr := addBucketPolicy(service, bucketName)

	if policyErr != nil {
		return policyErr
	}

	return nil
}

func addBucketPolicy(service *s3.S3, bucketName string) error {
	bucketACL :=
		`{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "AddPerm",
                "Effect": "Allow",
                "Principal": "*",
                "Action": "s3:GetObject",
                "Resource": "arn:aws:s3:::%s/*"
            }
        ]
    }`
	policy := fmt.Sprintf(bucketACL, bucketName)
	_, perr := service.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: &bucketName,
		Policy: &policy,
	})

	if perr != nil {
		fmt.Printf("Unable to set up ACL policy for %s: %+v\n", bucketName, perr)
	}
	return perr
}
