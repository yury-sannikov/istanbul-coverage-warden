package main

import (
	"fmt"

	"github.com/yury-sannikov/istanbul-coverage-warden/s3tools"
)

// CopyCommand CopyCommand
type CopyCommand struct {
	Bucket       string `short:"d" long:"bucket" description:"Destination S3 bucket to store converage information" required:"true"`
	Branch       string `short:"b" long:"branch" description:"Git branch name with coverage information" required:"true"`
	CoverageFile string `short:"f" long:"file" description:"Coverage file path" required:"true"`
	Region       string `short:"r" long:"region" description:"AWS S3 Region" default:"us-east-1"`
}

var copyCommand CopyCommand

// Execute copy code coverate to S3 bucket
func (opts *CopyCommand) Execute(args []string) error {

	connectionOptions := &s3tools.S3ConnectionOptions{Region: &opts.Region}

	sess, err := s3tools.CreateS3Session(connectionOptions)

	if err != nil {
		panic(err)
	}

	if err := s3tools.CheckOrCreateBucket(sess, opts.Bucket); err != nil {
		panic(err)
	}
	bucketPath := fmt.Sprintf("%s/%s", opts.Bucket, opts.Branch)
	fmt.Printf("Copying coverage information into %s bucket...\n", bucketPath)

	uploadKey := fmt.Sprintf("%s/coverage.xml", opts.Branch)
	if err := s3tools.UploadFile(sess, opts.Bucket, uploadKey, opts.CoverageFile); err != nil {
		panic(err)
	}
	return nil
}

func init() {
	parser.AddCommand("copy",
		"Copy code coverage to S3 bucket",
		"The copy command add or replace code coverage file for a specified branch to the S3 Bucket. Use -d to specify destination bucket and -b to specify branch name.",
		&copyCommand)
}
