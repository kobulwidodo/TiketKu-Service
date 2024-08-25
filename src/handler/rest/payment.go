package rest

import (
	"go-clean/src/business/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Payment
// @Description Create New Payment
// @Security BearerAuth
// @Tags Payment
// @Param payment body entity.CreatePaymentParam true "payment info"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.PaymentRes{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/pay [POST]
func (r *rest) CreatePayment(ctx *gin.Context) {
	var paymentParam entity.CreatePaymentParam
	if err := ctx.ShouldBindJSON(&paymentParam); err != nil {
		r.httpRespError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	payment, err := r.uc.Payment.Create(ctx.Request.Context(), paymentParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "sucessfully create a payment", payment)
}
