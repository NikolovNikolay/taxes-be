package statements

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
	"taxes-be/internal/atleastonce"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/models"
)

const processStatementKey = "process_statement"

type StatementManager struct {
	inquiryStore inquiriesdao.Store
	s3           *s3.S3
	alo          atleastonce.Doer
	s3BucketName string
}

func NewStatementManager(
	alo atleastonce.Doer,
	s3 *s3.S3,
	s3BucketName string,
) *StatementManager {
	sm := &StatementManager{
		alo:          alo,
		s3BucketName: s3BucketName,
		s3:           s3,
	}
	alo.RegisterHandler(processStatementKey, sm.handleProcessStatement)
	return sm
}

func (sm *StatementManager) handleProcessStatement(ctx context.Context, uuid uuid.UUID) error {
	inq, err := sm.inquiryStore.FindInquiry(ctx, uuid)
	if err != nil {
		return nil
	}

	// TODO: init statements processing logic

	fn := strings.Split(inq.Files, ",")
	for i := range fn {
		logrus.Info(fmt.Sprintf("processing file: %s", fn[i]))
	}

	err = sm.deleteFromS3(inq)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StatementManager) deleteFromS3(inq *models.Inquiry) error {
	iter := s3manager.NewDeleteListIterator(sm.s3, &s3.ListObjectsInput{
		Prefix: &inq.Prefix,
		Bucket: aws.String(sm.s3BucketName),
	})

	err := s3manager.NewBatchDeleteWithClient(sm.s3).Delete(aws.BackgroundContext(), iter)
	if err != nil {
		return err
	}

	return nil
}
