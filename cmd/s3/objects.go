package s3

import (
	"cptdom/sasha/utils"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getAllFiles(svc *s3.S3, b *string, prefix *string, delim *string) ([]*s3.ListObjectsV2Output, error) {
	defaultPrefix := ""
	if *prefix == "/" {
		prefix = &defaultPrefix
	}
	var result []*s3.ListObjectsV2Output
	params := &s3.ListObjectsV2Input{
		Bucket:    aws.String(*b),
		Delimiter: delim,
		Prefix:    prefix,
	}
	err := svc.ListObjectsV2Pages(params,
		func(page *s3.ListObjectsV2Output, empty bool) bool {
			if len(page.Contents) == 0 {
				return false
			}
			result = append(result, page)
			return true
		})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getAllDirs(svc *s3.S3, b *string, prefix *string, delim *string) (*s3.ListObjectsV2Output, error) {
	defaultPrefix := ""
	if *prefix == "/" {
		prefix = &defaultPrefix
	}
	result, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:    aws.String(*b),
		Delimiter: delim,
		Prefix:    prefix,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getObject(svc *s3.S3, b *string, filename *string) (*s3.GetObjectOutput, error) {
	result, err := svc.GetObject(&s3.GetObjectInput{
		Key:    filename,
		Bucket: b,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DescribeObject(svc *s3.S3, b *string, filename string) error {
	result, err := getObject(svc, b, &filename)
	splitName := strings.Split(filename, "/")
	if err != nil {
		fmt.Printf("%s: No such file.\n", splitName[len(splitName)-1])
		return err
	}
	fmt.Println(utils.CreateLine())
	fmt.Printf("Name: %v\n", filename)
	fmt.Printf("ETag: %s\n", *result.ETag)
	fmt.Printf("Type: %s\n", *result.ContentType)
	fmt.Printf("Last modified: %s\n", *result.LastModified)
	fmt.Println("Metadata: ")
	if len(result.Metadata) != 0 {
		for key, el := range result.Metadata {
			fmt.Printf("%v: %v\n", key, *el)
		}
	} else {
		fmt.Print("<empty>\n")
	}

	fmt.Println(utils.CreateLine())
	return nil
}

func ListDirs(svc *s3.S3, b *string, prefix *string) (map[string]struct{}, error) {
	dirDelim := "/"
	result, err := getAllDirs(svc, b, prefix, &dirDelim)
	if err != nil {
		fmt.Printf("Got an error retrieving directories: %v\n", err)
		return nil, err
	}
	results := make(map[string]struct{})
	var dummyStruct struct{}
	for _, r := range result.CommonPrefixes {
		results[*r.Prefix] = dummyStruct
	}
	return results, nil
}

func ListFiles(svc *s3.S3, b *string, prefix *string) {
	var fileDelim string
	result, err := getAllFiles(svc, b, prefix, &fileDelim)
	if err != nil {
		fmt.Printf("Got an error retrieving files: %v\n", err)
		return
	}
	var keycount int64 = 0
	if *prefix == "/" {
		*prefix = ""
	}
	for _, page := range result {
		for _, f := range page.Contents {
			splitF := strings.Split(*f.Key, "/")
			prefixLength := len(strings.Split(*prefix, "/"))
			if len(splitF) == prefixLength {
				fmt.Printf("%-50s %s\n", splitF[prefixLength-1], f.LastModified.Format(utils.DateFormat))
			}
		}
		keycount += *page.KeyCount
	}
	fmt.Printf("\n%v files in total.\n", keycount)
}
