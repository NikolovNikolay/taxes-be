package statements

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
	"taxes-be/internal/atleastonce"
	awsutil "taxes-be/internal/aws"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/sendgrid"
)

const processStatementKey = "process_statement"

type StatementManager struct {
	inquiryStore *inquiriesdao.Store
	couponStore  *couponsdao.Store
	s3Manager    *awsutil.S3Manager
	alo          atleastonce.Doer
	s3BucketName string
	mailer       *sendgrid.Mailer
}

func NewStatementManager(
	alo atleastonce.Doer,
	s3Manager *awsutil.S3Manager,
	s3BucketName string,
	mailer *sendgrid.Mailer,
	inquiryStore *inquiriesdao.Store,
	couponStore *couponsdao.Store,
) *StatementManager {
	sm := &StatementManager{
		couponStore:  couponStore,
		inquiryStore: inquiryStore,
		alo:          alo,
		s3BucketName: s3BucketName,
		s3Manager:    s3Manager,
		mailer:       mailer,
	}
	alo.RegisterHandler(processStatementKey, sm.handleProcessStatement)
	return sm
}

func (sm *StatementManager) handleProcessStatement(ctx context.Context, id uuid.UUID) error {
	inq, err := sm.inquiryStore.FindInquiry(ctx, id)
	if err != nil {
		return err
	}

	if !inq.Paid {
		return fmt.Errorf("inquiry not paid")
	}

	var couponID string
	if !inq.GeneratedWithCoupon {
		c, err := sm.couponStore.FindCouponByRequestID(ctx, id)
		if err != nil {
			return err
		}
		couponID = c.ID
	}

	fn := strings.Split(inq.Files, ",")
	rp := NewReportProcessor(inq.Year, inq.ID)

	for i := range fn {
		err = sm.handleSingleFile(fn[i], inq.Type, rp)
		if err != nil {
			return err
		}
	}

	err = rp.CalculateTaxes()
	if err != nil {
		return err
	}

	err = sm.mailer.SendReportMail(inq.Year, inq.FullName, inq.Email, rp.report, couponID)
	if err != nil {
		return err
	}

	err = sm.s3Manager.BatchDelete(sm.s3BucketName, inq.Prefix)
	if err != nil {
		return err
	}

	return nil
}

func (sm *StatementManager) handleSingleFile(
	fileName string,
	sType int,
	rp *ReportProcessor,
) error {
	_, err := sm.s3Manager.DownloadFile(sm.s3BucketName, fileName)
	if err != nil {
		return err
	}

	err = rp.ParseLines(fileName, sType)
	if err != nil {
		return err
	}

	_ = os.Remove(fileName)
	return nil
}
