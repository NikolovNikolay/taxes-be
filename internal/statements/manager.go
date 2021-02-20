package statements

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"taxes-be/internal/atleastonce"
	awsutil "taxes-be/internal/aws"
	"taxes-be/internal/core"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/models"
	"taxes-be/internal/sendgrid"
	"time"
)

const processStatementKey = "process_statement"
const deleteStatementsKey = "delete_statements"

type StatementManager struct {
	inquiryStore *inquiriesdao.Store
	couponStore  *couponsdao.Store
	aloStore     atleastonce.Store
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
	aloStore atleastonce.Store,
) *StatementManager {
	sm := &StatementManager{
		couponStore:  couponStore,
		inquiryStore: inquiryStore,
		alo:          alo,
		s3BucketName: s3BucketName,
		s3Manager:    s3Manager,
		mailer:       mailer,
		aloStore:     aloStore,
	}

	alo.RegisterHandler(processStatementKey, sm.handleProcessStatement)
	alo.RegisterHandler(deleteStatementsKey, sm.handleDeleteStatements)
	return sm
}

func (sm *StatementManager) handleDeleteStatements(ctx context.Context, id uuid.UUID) error {
	logrus.
		WithContext(ctx).
		WithField("requestID", id)

	inq, err := sm.inquiryStore.FindInquiry(ctx, id)
	if err != nil {
		logrus.
			WithError(err).
			Error("error while attempting to fetch inquiry from store")
		return err
	}

	if time.Now().UnixNano()-inq.ModifiedAt.UnixNano() < (48 * time.Hour).Nanoseconds() {
		return core.AsNoLogError(fmt.Errorf("can't delete statements yet"))
	}

	return sm.deleteInquiryFiles(inq)
}

func (sm *StatementManager) handleProcessStatement(ctx context.Context, id uuid.UUID) error {
	logrus.
		WithContext(ctx).
		WithField("requestID", id)

	inq, err := sm.inquiryStore.FindInquiry(ctx, id)
	if err != nil {
		logrus.
			WithError(err).
			Error("error while attempting to fetch inquiry from store")
		return err
	}

	if !inq.Paid {
		if time.Now().UnixNano()-inq.ModifiedAt.UnixNano() < (48 * time.Hour).Nanoseconds() {
			err = sm.aloStore.Save(ctx, atleastonce.Task{
				Key:  deleteStatementsKey,
				ID:   uuid.MustParse(inq.ID),
				Done: false,
			})
			if err != nil {
				logrus.WithError(err).Error("error when attempting to add delete statements task")
				return err
			}
		}
		return core.AsNoLogError(fmt.Errorf("inquiry not paid"))
	}

	var couponID string
	if !inq.GeneratedWithCoupon {
		c, err := sm.couponStore.FindCouponByRequestID(ctx, id)
		if err != nil {
			logrus.WithError(err).Error("error while calling coupon store")
			return err
		}
		couponID = c.ID
	}

	fn := strings.Split(inq.Files, ",")
	rp := NewReportProcessor(inq.Year, inq.ID)

	for i := range fn {
		err = sm.handleSingleFile(fn[i], inq.Type, rp)
		if err != nil {
			logrus.WithError(err).Error("error while getting statements to parse")
			return err
		}
	}

	err = rp.CalculateTaxes()
	if err != nil {
		logrus.WithError(err).Error("error while calculating taxes")
		return err
	}

	err = sm.aloStore.Save(ctx, atleastonce.Task{
		Key:  deleteStatementsKey,
		ID:   uuid.MustParse(inq.ID),
		Done: false,
	})
	if err != nil {
		logrus.WithError(err).Error("error when attempting to add delete statements task")
		return err
	}

	err = sm.mailer.SendReportMail(inq.Year, inq.FullName, inq.Email, rp.report, couponID)
	if err != nil {
		logrus.Error("error while calculating taxes")
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

func (sm *StatementManager) deleteInquiryFiles(inq *models.Inquiry) error {
	err := sm.s3Manager.BatchDelete(sm.s3BucketName, inq.Prefix)
	if err != nil {
		logrus.WithError(err).Error("error while attempting to delete statements after calculated taxes")
		return err
	}

	return nil
}
