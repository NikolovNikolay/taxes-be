package statementsview

import (
	"github.com/labstack/echo"
	"taxes-be/internal/atleastonce"
	awsutil "taxes-be/internal/aws"
	"taxes-be/internal/inquiries/inquiriesdao"
)

func RegisterStatementsEndpoints(
	group *echo.Group,
	s3Manager *awsutil.S3Manager,
	s3BucketName string,
	inquiryStore *inquiriesdao.Store,
	aloStore atleastonce.Store,
) {
	ue := NewUploadFilesEndpoint(inquiryStore, aloStore, s3Manager, s3BucketName)
	group.POST("/upload", ue.ServeHTTP)
}
