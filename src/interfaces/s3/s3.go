package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client struct {
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func New() (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		return nil, err
	}
	return &S3Client{
		uploader:   s3manager.NewUploader(sess),
		downloader: s3manager.NewDownloader(sess),
	}, nil
}

func (s *S3Client) Download(bucket, key string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := s.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
