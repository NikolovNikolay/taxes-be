package paymentsview

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
	"taxes-be/internal/core"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/inquiries/inquiriesdao"
	"taxes-be/internal/models"
	util "taxes-be/utils"
)

type ObtainCouponEndpoint struct {
	couponStore  *couponsdao.Store
	inquiryStore *inquiriesdao.Store
}

func NewObtainCouponEndpoint(
	couponStore *couponsdao.Store,
	inquiryStore *inquiriesdao.Store,
) *ObtainCouponEndpoint {
	return &ObtainCouponEndpoint{
		couponStore:  couponStore,
		inquiryStore: inquiryStore,
	}
}

type ObtainCouponRequest struct {
	RequestID uuid.UUID `json:"request_id" validate:"required"`
}

type ObtainCouponResponse struct {
	Coupon string `json:"coupon"`
}

func (ep *ObtainCouponEndpoint) ServeHTTP(c echo.Context) error {
	r := &ObtainCouponRequest{}
	if err := util.ValidateRequest(c, r); err != nil {
		return core.CtxAware(c.Request().Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Internal: err,
			Message:  "could not parse json body",
		})
	}

	inquiry, err := ep.inquiryStore.FindInquiry(c.Request().Context(), r.RequestID)
	if err != nil {
		if core.IsNotFound(err) {
			return core.CtxAware(c.Request().Context(), &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Internal: err,
				Message:  "invalid request id",
			})
		}

		return core.CtxAware(c.Request().Context(), &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Internal: err,
			Message:  "could not obtain coupon",
		})
	}

	if inquiry.GeneratedWithCoupon {
		return c.JSON(http.StatusOK, &ObtainCouponResponse{
			Coupon: "",
		})
	}

	coupon, err := ep.couponStore.FindCouponByRequestID(c.Request().Context(), r.RequestID)
	if err != nil {
		if !core.IsNotFound(err) {
			return core.CtxAware(c.Request().Context(), &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Internal: err,
				Message:  "could not obtain coupon",
			})
		}
	}

	if coupon == nil {
		couponID := uuid.New()
		err = ep.inquiryStore.InTransaction(c.Request().Context(), func(ctx context.Context) error {
			err = ep.couponStore.UpsertCoupon(ctx, &models.Coupon{
				ID:              couponID.String(),
				MaxAttempts:     3,
				Attempts:        1,
				Type:            inquiry.Type,
				ParentRequestID: r.RequestID.String(),
				Email:           inquiry.Email,
			})

			if err != nil {
				return err
			}

			inquiry.Paid = true
			err = ep.inquiryStore.UpsertInquiry(ctx, inquiry)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return core.CtxAware(c.Request().Context(), &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Internal: err,
				Message:  "could not obtain coupon",
			})
		}
		return c.JSON(http.StatusCreated, &ObtainCouponResponse{
			Coupon: couponID.String(),
		})
	} else {
		return core.CtxAware(c.Request().Context(), &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Internal: err,
			Message:  "coupon already exists for provided request",
		})
	}
}
