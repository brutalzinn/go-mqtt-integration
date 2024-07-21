package awshelper

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/samber/lo"
)

func S3PutObject(ctx context.Context, bucketRegion string, bucket string, fileName string, data []byte) error {
	s3Client := clientConfig(bucketRegion)
	fileType := http.DetectContentType(data)
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &fileName,
		Body:        bytes.NewReader(data),
		ContentType: aws.String(fileType),
	})
	return err
}

func S3ListObject(ctx context.Context, bucketRegion string, bucket string) ([]string, error) {
	s3Client := clientConfig(bucketRegion)
	list, err := s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, err
	}
	ret := lo.Map(list.Contents, func(item types.Object, index int) string {
		return *item.Key
	})
	return ret, err
}

func S3RenameObject(ctx context.Context, bucketRegion string, bucket string, fileName string, data []byte) error {
	return nil
}

func S3DeleteObject(ctx context.Context, bucketRegion string, bucket string, fileName string) error {
	s3Client := clientConfig(bucketRegion)
	_, err := s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
	})
	return err
}

func S3GetObject(ctx context.Context, bucketRegion string, bucket string, fileName string) ([]byte, error) {
	s3Client := clientConfig(bucketRegion)
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &fileName,
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
