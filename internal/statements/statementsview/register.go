package statementsview

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/labstack/echo"
)

func RegisterStatementsEndpoints(group *echo.Group, session *session.Session, s3BucketName string) {
	ue := NewUploadFilesEndpoint(session, s3BucketName)
	group.POST("/upload", ue.ServeHTTP)
}
