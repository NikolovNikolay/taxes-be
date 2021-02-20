package paymentsview

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/checkout/session"
	"net/http"
	"taxes-be/internal/core"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/inquiries/inquiriesdao"
	util "taxes-be/utils"
)

type CreateCheckoutSessionEndpoint struct {
	websiteBaseUrl string
	couponStore    *couponsdao.Store
	inquiryStore   *inquiriesdao.Store
}

func NewCreateCheckoutSessionEndpoint(
	websiteBaseUrl string,
	inquiryStore *inquiriesdao.Store,
	couponStore *couponsdao.Store,
) *CreateCheckoutSessionEndpoint {
	return &CreateCheckoutSessionEndpoint{
		websiteBaseUrl: websiteBaseUrl,
		inquiryStore:   inquiryStore,
		couponStore:    couponStore,
	}
}

type CreateCheckoutSessionRequest struct {
	RequestID uuid.UUID `json:"request_id" validate:"required"`
}

type CreateCheckoutSessionResponse struct {
	SessionID string `json:"id"`
}

func (ep *CreateCheckoutSessionEndpoint) ServeHTTP(c echo.Context) error {
	r := &CreateCheckoutSessionRequest{}
	if err := util.ValidateRequest(c, r); err != nil {
		return core.CtxAware(c.Request().Context(), &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Internal: err,
			Message:  "could not parse json body",
		})
	}

	inq, err := ep.inquiryStore.FindInquiry(c.Request().Context(), r.RequestID)
	if err != nil {
		if core.IsNotFound(err) {
			return core.CtxAware(c.Request().Context(), &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Internal: err,
				Message:  "invalid request ID",
			})
		}

		return core.CtxAware(c.Request().Context(), &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Internal: err,
			Message:  "could not create checkout session",
		})
	}

	callbackTemplate := "%s/checkout?success=%v&request_id=%s"
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("bgn"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("statements report"),
					},
					UnitAmount: stripe.Int64(10 * 100),
				},
				Quantity: stripe.Int64(1),
			},
		},
		CustomerEmail: &inq.Email,
		SuccessURL: stripe.String(
			fmt.Sprintf(callbackTemplate, ep.websiteBaseUrl, true, r.RequestID),
		),
		CancelURL: stripe.String(
			fmt.Sprintf(callbackTemplate, ep.websiteBaseUrl, false, r.RequestID),
		),
	}

	sess, err := session.New(params)
	if err != nil {
		return err
	}

	data := CreateCheckoutSessionResponse{
		SessionID: sess.ID,
	}

	return c.JSON(http.StatusOK, data)
}
