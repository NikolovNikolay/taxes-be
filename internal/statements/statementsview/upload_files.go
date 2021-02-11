package statementsview

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"strings"
	"taxes-be/internal/atleastonce"
	"taxes-be/internal/core"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/models"
	"time"
)

type UploadFilesEndpoint struct {
	inquiryStore *inquiriesdao.Store
	aloStore     atleastonce.Store
	awsSession   *session.Session
	bucketName   string
}

func NewUploadFilesEndpoint(
	inquiryStore *inquiriesdao.Store,
	aloStore atleastonce.Store,
	awsSession *session.Session,
	bucketName string,
) *UploadFilesEndpoint {
	return &UploadFilesEndpoint{
		inquiryStore: inquiryStore,
		aloStore:     aloStore,
		bucketName:   bucketName,
		awsSession:   awsSession,
	}
}

func (ep *UploadFilesEndpoint) ServeHTTP(c echo.Context) error {
	req := c.Request()
	err := req.ParseMultipartForm(0)

	if err != nil {
		m := "error processing statement files"
		if err.Error() == "http: request body too large" {
			m = err.Error()
		}
		logrus.Errorf("error parsing statement files: %v", err)
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Internal: err,
			Message:  m,
		})
	}

	mf := req.MultipartForm
	files := mf.File["statements"]
	allFileNames := make([]string, 0)

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	for i := range files {
		t := readFileExtension(files[i])
		if t == nil || !isSupportedExtension(t) {
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "unsupported statement type",
			})
		}

		file, err := files[i].Open()
		if err != nil {
			logrus.Errorf("error opening file: %d", i)
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "error reading files",
			})
		}

		fileName := fmt.Sprintf("%s_%s.%s", prefix, uuid.New().String(), t.Extension)
		uploader := s3manager.NewUploader(ep.awsSession)

		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(ep.bucketName),
			ACL:    aws.String("bucket-owner-full-control"),
			Key:    aws.String(fileName),
			Body:   file,
		})

		if err != nil {
			logrus.Errorf("could not upload files: %v", err)
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "could not process statements",
			})
		}

		allFileNames = append(allFileNames, fileName)

		file.Close()
	}

	id := uuid.New()
	err = ep.inquiryStore.InTransaction(req.Context(), func(ctx context.Context) error {
		err = ep.inquiryStore.AddInquiry(ctx, &models.Inquiry{
			ID:     id.String(),
			UserID: uuid.New().String(),
			Prefix: prefix,
			Files:  strings.Join(allFileNames, ","),
		})

		if err != nil {
			return err
		}

		err = ep.aloStore.Save(ctx, atleastonce.Task{
			Key: "process_statement",
			ID:  id,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logrus.Errorf("error processing statements: %v", err)
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "could not process statements",
		})
	}

	return c.NoContent(http.StatusCreated)
}

func readFileExtension(f *multipart.FileHeader) *types.Type {
	file, err := f.Open()
	defer file.Close()

	if err != nil || file == nil {
		logrus.Errorf("error parsing file type: %s", f.Filename)
		return nil
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		logrus.Errorf("error reading file: %s", f.Filename)
		return nil
	}
	kind, _ := filetype.Match(buff)
	if kind == filetype.Unknown {
		return nil
	}

	return &kind
}

func isSupportedExtension(t *types.Type) bool {
	return t.Extension == "pdf" || t.Extension == "xlsx" || t.Extension == "xls"
}
