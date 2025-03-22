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
	// S3アップロード
	S3Upload(key string, fileHeader *multipart.FileHeader) error
	// S3ダウンロード
	S3Download(fileName string) ([]byte, error)
	// S3複数ファイル削除
	S3DeleteFiles(keys []string) error
}

type AWSRepository struct{}

func NewAWSRepository() IAWSRepository {
	return &AWSRepository{}
}

// S3アップロード
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

// S3ダウンロード
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

// S3から複数のファイルを削除
func (a *AWSRepository) S3DeleteFiles(keys []string) error {
	// AWSセッションを作成
	sess, sessErr := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
		Region:      aws.String(os.Getenv("AWS_REGION")),
	})
	if sessErr != nil {
		log.Printf("%v", sessErr)
		return sessErr
	}

	// S3サービスクライアントを作成
	svc := s3.New(sess)

	// 削除対象のオブジェクトリストを作成
	var objects []*s3.ObjectIdentifier
	for _, key := range keys {
		objects = append(objects, &s3.ObjectIdentifier{Key: aws.String(key)})
	}

	// 複数ファイルをS3から削除
	_, delErr := svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Delete: &s3.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	if delErr != nil {
		log.Printf("%v", delErr)
		return delErr
	}

	return nil
}
