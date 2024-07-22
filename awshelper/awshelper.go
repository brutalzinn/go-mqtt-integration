package awshelper

import (
	"bytes"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/samber/lo"
)

func S3PutObject(bucketRegion string, bucket string, fileName string, data []byte) error {
	s3Client := clientConfig(bucketRegion)
	fileType := http.DetectContentType(data)
	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(fileType),
	})
	return err
}

func S3ListObject(bucketRegion string, bucket string) ([]string, error) {
	s3Client := clientConfig(bucketRegion)
	list, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}
	ret := lo.Map(list.Contents, func(item *s3.Object, index int) string {
		return *item.Key
	})
	return ret, err
}

func S3DeleteObject(bucketRegion string, bucket string, fileName string) error {
	s3Client := clientConfig(bucketRegion)
	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	return err
}

func S3GetObject(bucketRegion string, bucket string, fileName string) ([]byte, error) {
	s3Client := clientConfig(bucketRegion)
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()
	objBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}
	return objBytes, nil
}
