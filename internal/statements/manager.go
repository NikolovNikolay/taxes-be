package statements

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"taxes-be/internal/atleastonce"
	awsutil "taxes-be/internal/aws"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/sendgrid"
)

const processStatementKey = "process_statement"

//const (
//	pdf                 = "pdf"
//	xlsx                = "xlsx"
//	xls                 = "xls"
//	processStatementKey = "process_statement"
//)
//
//var (
//	excelReader      = reader.NewExcelReader()
//	pdfReader        = reader.NewPDFReader()
//	supportedFormats = map[string]reader.Reader{
//		pdf:  pdfReader,
//		xlsx: excelReader,
//		xls:  excelReader,
//	}
//)

type StatementManager struct {
	inquiryStore inquiriesdao.Store
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
) *StatementManager {
	sm := &StatementManager{
		alo:          alo,
		s3BucketName: s3BucketName,
		s3Manager:    s3Manager,
		mailer:       mailer,
	}
	alo.RegisterHandler(processStatementKey, sm.handleProcessStatement)
	return sm
}

func (sm *StatementManager) handleProcessStatement(ctx context.Context, uuid uuid.UUID) error {
	inq, err := sm.inquiryStore.FindInquiry(ctx, uuid)
	if err != nil {
		return nil
	}

	fn := strings.Split(inq.Files, ",")
	rp := NewReportProcessor(2020, inq.ID)

	for i := range fn {
		err = sm.handleSingleFile(fn[i], inq.Type, rp)
		if err != nil {
			return err
		}
		logrus.Info("done")
	}

	err = rp.CalculateTaxes()
	if err != nil {
		return err
	}

	err = sm.mailer.SendReportMail(2020, "Nikolay Nikolov", "nikolayvnikolov@protonmail.com", rp.report)
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
