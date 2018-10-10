package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/yury-sannikov/istanbul-coverage-warden/s3tools"
	"github.com/yury-sannikov/istanbul-coverage-warden/xmltools"
)

// CheckCommand parameters
type CheckCommand struct {
	Bucket       string `short:"d" long:"bucket" description:"Source S3 bucket to get converage information" required:"true"`
	Branch       string `short:"b" long:"branch" description:"Git branch name with coverage information to compare" required:"true"`
	CoverageFile string `short:"f" long:"file" description:"Coverage file path" required:"true"`
	Region       string `short:"r" long:"region" description:"AWS S3 Region" default:"us-east-1"`
}

var checkCommand CheckCommand

type coberturaClassMap map[string]xmltools.CoberturaClass

const _NotFoundError = "NoSuchKey"

const _DropToleranceRate = -0.008

// Execute copy code coverate to S3 bucket
func (opts *CheckCommand) Execute(args []string) error {

	connectionOptions := &s3tools.S3ConnectionOptions{Region: &opts.Region}

	sess, err := s3tools.CreateS3Session(connectionOptions)

	if err != nil {
		panic(err)
	}

	tmpfile, err := ioutil.TempFile("", "istanbul-coverage-warden")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	coverageFile := fmt.Sprintf("%s/coverage.xml", opts.Branch)
	fmt.Printf("Downloading coverage file for branch %s from %s ... ", opts.Branch, opts.Bucket)

	size, err := s3tools.DownloadFile(sess, opts.Bucket, coverageFile, tmpfile.Name())

	if err != nil {
		fmt.Printf("\nUnable to download coverage file %s from bucket %s.\n", coverageFile, opts.Bucket)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == _NotFoundError {
				// Do not error if no coverage file found. It might not be enabled
				fmt.Println("No coverage file has been stored.")
				os.Exit(0)
			}
		}
		panic(err)
	}

	fmt.Printf("downloaded %d kb\n", size/1024)

	currentClasses, err := xmltools.BuildCoberturaClasses(opts.CoverageFile)
	if err != nil {
		fmt.Println("Error opening current cobertura report file:", err)
		os.Exit(2)
	}

	prevClasses, err := xmltools.BuildCoberturaClasses(tmpfile.Name())
	if err != nil {
		fmt.Println("Error opening previous cobertura report file:", err)
		os.Exit(3)
	}

	currentClassMap := make(coberturaClassMap)
	for _, item := range currentClasses {
		currentClassMap[item.FileName] = item
	}
	result := compare(prevClasses, currentClassMap)

	if !result {
		fmt.Println("Code coverage drop has been detected.")
		os.Exit(1)
	}
	fmt.Println("No coverage drop has been detected.")

	return nil
}

func init() {
	parser.AddCommand("check",
		"Get original coverage information from destination branch and check for coverage drops",
		"The check command fetch code coverage file for a specified branch from the S3 Bucket and check code coverage drops. Use -d to specify source bucket, -b to specify branch name and -f to point to the coverage for check.",
		&checkCommand)
}

func compare(previousCoberturaData []xmltools.CoberturaClass, newCoberturaMap coberturaClassMap) bool {
	var result = true
	for _, item := range previousCoberturaData {
		newItem, found := newCoberturaMap[item.FileName]
		if !found {
			fmt.Println("Unable to find coverage infrormation for ", item.FileName)
			continue
		}
		rateDiff := newItem.LineRate - item.LineRate
		if rateDiff < 0.0 {
			if _DropToleranceRate > rateDiff {
				fmt.Printf("Code coverage dropped for %s from %.2f%% to %.2f%%\n", item.FileName, item.LineRate*100.0, newItem.LineRate*100.0)
				result = false
			} else {
				fmt.Printf("Insignificant code coverage dropped for %s from %.2f%% to %.2f%%\n", item.FileName, item.LineRate*100.0, newItem.LineRate*100.0)
			}
		}
	}
	return result
}
