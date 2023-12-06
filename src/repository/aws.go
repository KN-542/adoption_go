package repository

import (
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type IAWSRepository interface {
	S3Upload(key string, fileHeader *multipart.FileHeader) error
	S3Download(fileName string) ([]byte, error)
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
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		Region:      aws.String(os.Getenv("AWS_REGION"))},
	)
	if err != nil {
		log.Printf("%v", err)
		return err
	}

	uploader := s3manager.NewUploader(s)
	_, err2 := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   file,
	})
	if err2 != nil {
		log.Printf("%v", err2)
		return err2
	}

	return nil
}

func (a *AWSRepository) S3Download(fileName string) ([]byte, error) {
	// AWSセッションを作成
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		Region:      aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	// S3サービスクライアントを作成
	svc := s3.New(sess)

	// S3からファイルを取得
	res, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		// エラーの処理
		log.Printf("%v", err)
		return nil, err
	}

	return data, nil
}
