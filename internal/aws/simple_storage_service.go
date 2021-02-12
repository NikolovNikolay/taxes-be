package awsutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
)

type S3Manager struct {
	session    *session.Session
	s3         *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3Manager(session *session.Session) *S3Manager {
	return &S3Manager{
		session:    session,
		s3:         s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

func (s *S3Manager) UploadMultipartFile(bucket, acl, fileName string, file io.Reader) error {
	_, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
		Key:    aws.String(fileName),
		Body:   file,
	})
	return err
}

func (s *S3Manager) BatchDelete(bucket, filesPrefix string) error {
	var fp *string = nil
	if filesPrefix != "" {
		fp = &filesPrefix
	}

	iter := s3manager.NewDeleteListIterator(s.s3, &s3.ListObjectsInput{
		Prefix: fp,
		Bucket: aws.String(bucket),
	})

	err := s3manager.NewBatchDeleteWithClient(s.s3).Delete(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Manager) DownloadFile(bucket, fileName string) (*os.File, error) {
	file, err := os.Create(fileName)

	if err != nil {
		return nil, fmt.Errorf("could not create temp file for: %s", fileName)
	}
	_, err = s.downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return nil, fmt.Errorf("could not download file: %s", fileName)
	}

	_ = file.Close()
	return file, nil
}
