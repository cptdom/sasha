package s3

import (
	"github.com/aws/aws-sdk-go/service/s3"

	"fmt"
)

func getAllBuckets(svc *s3.S3) (*s3.ListBucketsOutput, error) {

	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAllBuckets(svc *s3.S3) ([]*s3.Bucket, error) {
	res, err := getAllBuckets(svc)
	if err != nil {
		fmt.Printf("Got an error retrieving buckets: %v\n", err)
		return nil, err
	}
	var result []*s3.Bucket
	for _, bucket := range res.Buckets {
		result = append(result, bucket)
	}
	return result, nil
}
