package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get List Payment Option
// @Description Get List All Payment Option
// @Security BearerAuth
// @Tags PaymentOption
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.PaymentOption{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/payment-options [GET]
func (r *rest) GetListPaymentOption(ctx *gin.Context) {
	paymentOptions, err := r.uc.PaymentOption.GetList(entity.PaymentOptionParam{})
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get list all event", paymentOptions)
}
