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
	// S3アップロード*
	S3Upload(key string, fileHeader *multipart.FileHeader) error
	// S3ダウンロード*
	S3Download(fileName string) ([]byte, error)
}

type AWSRepository struct{}

func NewAWSRepository() IAWSRepository {
	return &AWSRepository{}
}

// S3 Upload
func (a *AWSRepository) S3Upload(key string, fileHeader *multipart.FileHeader) error {
	// ファイルを開く
	file, fileErr := fileHeader.Open()
	if fileErr != nil {
		log.Printf("%v", fileErr)
		return fileErr
	}
	defer file.Close()

	// AWSセッションを作成（東京リージョン）
	s, sErr := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		Region:      aws.String(os.Getenv("AWS_REGION"))},
	)
	if sErr != nil {
		log.Printf("%v", sErr)
		return sErr
	}

	uploader := s3manager.NewUploader(s)
	_, upErr := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(key),
		Body:   file,
	})
	if upErr != nil {
		log.Printf("%v", upErr)
		return upErr
	}

	return nil
}

func (a *AWSRepository) S3Download(fileName string) ([]byte, error) {
	// AWSセッションを作成
	sess, sessErr := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		Region:      aws.String(os.Getenv("AWS_REGION")),
	})
	if sessErr != nil {
		log.Printf("%v", sessErr)
		return nil, sessErr
	}

	// S3サービスクライアントを作成
	svc := s3.New(sess)

	// S3からファイルを取得
	res, resErr := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(fileName),
	})
	if resErr != nil {
		return nil, resErr
	}
	defer res.Body.Close()

	data, dataErr := io.ReadAll(res.Body)
	if dataErr != nil {
		// エラーの処理
		log.Printf("%v", dataErr)
		return nil, dataErr
	}

	return data, nil
}
