package statementsview

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/null/v8"
	"net/http"
	"strconv"
	"strings"
	"taxes-be/internal/atleastonce"
	awsutil "taxes-be/internal/aws"
	"taxes-be/internal/core"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/models"
	files2 "taxes-be/utils/files"
	"time"
)

type UploadFilesEndpoint struct {
	inquiryStore *inquiriesdao.Store
	couponStore  *couponsdao.Store
	aloStore     atleastonce.Store
	s3Manager    *awsutil.S3Manager
	bucketName   string
}

func NewUploadFilesEndpoint(
	inquiryStore *inquiriesdao.Store,
	couponStore *couponsdao.Store,
	aloStore atleastonce.Store,
	s3Manager *awsutil.S3Manager,
	bucketName string,
) *UploadFilesEndpoint {
	return &UploadFilesEndpoint{
		inquiryStore: inquiryStore,
		aloStore:     aloStore,
		bucketName:   bucketName,
		s3Manager:    s3Manager,
		couponStore:  couponStore,
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
	sType := strings.Trim(req.FormValue("type"), " ")
	sYear := strings.Trim(req.FormValue("year"), " ")
	sEmail := strings.Trim(req.FormValue("email"), " ")
	sFullName := strings.Trim(req.FormValue("fullName"), " ")
	sCoupon := strings.Trim(req.FormValue("coupon"), " ")

	if len(files) == 0 {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "no attached statements",
			Internal: fmt.Errorf("no attached statements"),
		})
	}

	if sEmail == "" {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "missing email",
			Internal: fmt.Errorf("missing email"),
		})
	}

	year, err := strconv.Atoi(sYear)
	if err != nil {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "invalid year: " + sYear,
			Internal: fmt.Errorf("invalid year: " + sYear),
		})
	}

	inquiryType := int(mapType(sType))
	shouldNotPay := false
	var coupon *models.Coupon

	if sCoupon == "" {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "missing coupon code",
			Internal: fmt.Errorf("missing coupon code"),
		})
	}

	id, err := uuid.Parse(sCoupon)
	if err != nil {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "invalid coupon code",
			Internal: fmt.Errorf("invalid coupon code"),
		})
	}

	coupon, err = ep.couponStore.FindCoupon(req.Context(), id)
	if err != nil {
		if core.IsNotFound(err) {
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "invalid coupon code",
			})
		}
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "error processing request",
		})
	}

	if coupon.Attempts == coupon.MaxAttempts {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "expired coupon code",
		})
	}
	if inquiryType != coupon.Type {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "coupon code is not for this report type",
		})
	}

	if coupon.Email != sEmail {
		return core.CtxAware(req.Context(), &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "coupon code personalised for other user",
		})
	}
	if coupon.Attempts+1 <= coupon.MaxAttempts {
		shouldNotPay = true
	}

	allFileNames := make([]string, 0)
	prefix := fmt.Sprintf("%d", time.Now().Unix())
	for i := range files {
		ext, err := files2.GetMultipartFileExtension(files[i])
		if err != nil || !isSupportedExtension(ext, sType) {
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "unsupported statement type - PDF only for Revolut and Excel formats for eToro",
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

			_ = ep.s3Manager.BatchDelete(ep.bucketName, prefix)
			return core.CtxAware(req.Context(), &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "could not process statements",
			})
		}

		allFileNames = append(allFileNames, fileName)
		_ = file.Close()
	}

	inquiryID := uuid.New()
	err = ep.inquiryStore.InTransaction(req.Context(), func(ctx context.Context) error {
		err = ep.inquiryStore.UpsertInquiry(ctx, &models.Inquiry{
			ID:                  inquiryID.String(),
			UserID:              uuid.New().String(),
			Prefix:              prefix,
			Files:               strings.Join(allFileNames, ","),
			Type:                inquiryType,
			Year:                year,
			Email:               sEmail,
			FullName:            null.StringFrom(sFullName),
			Paid:                shouldNotPay,
			GeneratedWithCoupon: shouldNotPay,
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

		if shouldNotPay {
			coupon.Attempts = coupon.Attempts + 1
			err = ep.couponStore.UpsertCoupon(ctx, coupon)
			if err != nil {
				return err
			}
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
