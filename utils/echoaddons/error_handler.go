package echoaddons

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"taxes-be/internal/core"
	"taxes-be/internal/responses"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	ctx := c.Request().Context()
	// add detailed logs from the context
	if cae, ok := err.(*core.ContextAwareError); ok {
		ctx = cae.Ctx
	}

	respErr := asErrorResponse(ctx, err)

	le := logrus.WithContext(ctx).WithError(err)
	if respErr.Code >= 500 {
		errMsg := "operation failed"

		le.Error(errMsg)
	} else {
		le.Debug("request failed")
	}
	if err := c.JSON(respErr.Code, respErr); err != nil {
		logrus.WithContext(c.Request().Context()).WithError(err).Warn("writing error response failed")
	}
}

func asErrorResponse(ctx context.Context, err error) *responses.ErrorResponse {
	// unpack context aware errors
	if cae, ok := err.(*core.ContextAwareError); ok {
		ctx = cae.Ctx
		err = cae.Wrapped
	}

	if he, ok := err.(*echo.HTTPError); ok {
		if he.Code >= 500 {
			return &responses.ErrorResponse{
				Code:     http.StatusInternalServerError,
				Error:    fmt.Sprintf("Internal Server Error"),
				Internal: err,
			}
		}
		er := &responses.ErrorResponse{
			Code:        he.Code,
			Description: he.Message.(string),
			Internal:    he.Internal,
		}
		if he.Internal != nil {
			er.Error = he.Internal.Error()
		}
		return er
	}

	if core.IsValidationError(err) {
		return &responses.ErrorResponse{
			Code:        http.StatusBadRequest,
			Error:       fmt.Sprintf("Bad Request"),
			Description: err.Error(),
			Internal:    err,
		}
	}

	if core.IsNotFound(err) {
		return &responses.ErrorResponse{
			Code:        http.StatusNotFound,
			Error:       fmt.Sprintf("Not found"),
			Description: err.Error(),
			Internal:    err,
		}
	}

	return &responses.ErrorResponse{
		Code:     http.StatusInternalServerError,
		Error:    fmt.Sprintf("Internal Server Error"),
		Internal: err,
	}
}
