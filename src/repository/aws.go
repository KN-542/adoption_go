package repository

import (
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IAWSRepository interface {
	S3Upload(key string, fileHeader *multipart.FileHeader) error
}

type AWSRepository struct{}

func NewAWSRepository() IAWSRepository {
	return &AWSRepository{}
}

// S3 Upload
func (a *AWSRepository) S3Upload(key string, fileHeader *multipart.FileHeader) error {
	// ファイルを開く
	file, err := fileHeader.Open()
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	defer file.Close()

	// AWSセッションを作成（東京リージョン）
	s, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("<REDACTED>", "<REDACTED>/", ""),
		Region:      aws.String("ap-northeast-1")},
	)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	uploader := s3manager.NewUploader(s)
	_, err2 := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("adoption-resume"),
		Key:    aws.String(key),
		Body:   file,
	})
	if err2 != nil {
		log.Printf("%v", err2)
		return err2
	}

	return nil
}
