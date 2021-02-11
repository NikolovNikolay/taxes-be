package statementsview

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo"
	"taxes-be/internal/atleastonce"
	"taxes-be/internal/inquiries/inquiriesdao"
)

func RegisterStatementsEndpoints(
	group *echo.Group,
	session *session.Session,
	s3BucketName string,
	inquiryStore *inquiriesdao.Store,
	aloStore atleastonce.Store,
) {
	ue := NewUploadFilesEndpoint(inquiryStore, aloStore, session, s3BucketName)
	group.POST("/upload", ue.ServeHTTP)
}
