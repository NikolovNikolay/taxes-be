package paymentsview

import (
	"github.com/labstack/echo"
	"taxes-be/internal/coupons/couponsdao"
	"taxes-be/internal/inquiries/inquiriesdao"
)

func RegisterEndpoints(
	group *echo.Group,
	websiteBaseUrl string,
	couponStore *couponsdao.Store,
	inquiryStore *inquiriesdao.Store,
) {
	ccs := NewCreateCheckoutSessionEndpoint(websiteBaseUrl, couponStore)
	group.POST("/create-checkout-session", ccs.ServeHTTP)

	oce := NewObtainCouponEndpoint(couponStore, inquiryStore)
	group.POST("/obtain-coupon", oce.ServeHTTP)
}
