package statementsview

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"taxes-be/internal/atleastonce"
	awsutil "taxes-be/internal/aws"
	"taxes-be/internal/core"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/models"
	files2 "taxes-be/utils/files"
	"time"
)

type UploadFilesEndpoint struct {
	inquiryStore *inquiriesdao.Store
	aloStore     atleastonce.Store
	s3Manager    *awsutil.S3Manager
	bucketName   string
}

func NewUploadFilesEndpoint(
	inquiryStore *inquiriesdao.Store,
	aloStore atleastonce.Store,
	s3Manager *awsutil.S3Manager,
	bucketName string,
) *UploadFilesEndpoint {
	return &UploadFilesEndpoint{
		inquiryStore: inquiryStore,
		aloStore:     aloStore,
		bucketName:   bucketName,
		s3Manager:    s3Manager,
	}
}

type UploadResponse struct {
	RequestID uuid.UUID `json:"request_id"`
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
	sType := req.FormValue("type")

	if len(files) == 0 {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "no attached statements",
			Internal: fmt.Errorf("no attached statements"),
		})
	}

	allFileNames := make([]string, 0)

	prefix := fmt.Sprintf("%d", time.Now().Unix())
	for i := range files {
		ext, err := files2.GetMultipartFileExtension(files[i])
		if err != nil || !isSupportedExtension(ext, sType) {
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "unsupported statement type",
				Internal: err,
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

		fileName := fmt.Sprintf("%s_%s.%s", prefix, uuid.New().String(), ext)
		err = ep.s3Manager.UploadMultipartFile(ep.bucketName, "bucket-owner-full-control", fileName, file)
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

	inquiryID := uuid.New()
	err = ep.inquiryStore.InTransaction(req.Context(), func(ctx context.Context) error {
		err = ep.inquiryStore.AddInquiry(ctx, &models.Inquiry{
			ID:     inquiryID.String(),
			UserID: uuid.New().String(), // TODO: use real ID
			Prefix: prefix,
			Files:  strings.Join(allFileNames, ","),
			Type:   int(mapType(sType)),
		})

		if err != nil {
			return err
		}

		err = ep.aloStore.Save(ctx, atleastonce.Task{
			Key: "process_statement",
			ID:  inquiryID,
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

	return c.JSON(http.StatusAccepted, &UploadResponse{
		RequestID: inquiryID,
	})
}

func isSupportedExtension(ext, sType string) bool {
	if strings.ToLower(sType) == "etoro" {
		return isSupportedEToroExtension(ext)
	} else if strings.ToLower(sType) == "revolut" {
		return isSupportedRevolutExtension(ext)
	}
	return false
}

func isSupportedRevolutExtension(ext string) bool {
	return ext == "pdf"
}

func isSupportedEToroExtension(ext string) bool {
	return ext == "xlsx" || ext == "xls"
}

func mapType(sType string) core.StatementType {
	if strings.ToLower(sType) == "revolut" {
		return core.Revolut
	} else if strings.ToLower(sType) == "etoro" {
		return core.EToro
	} else {
		return core.Unknown
	}
}
