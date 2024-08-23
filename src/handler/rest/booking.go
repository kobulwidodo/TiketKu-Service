package rest

import (
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Create Booking
// @Description Create New Booking
// @Security BearerAuth
// @Tags Booking
// @Param event_id path integer true "event id"
// @Param category_id path integer true "category id"
// @Param booking body entity.CreateBookingParam true "booking info"
// @Produce json
// @Success 200 {object} entity.Response{data=entity.BookingResponse{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/event/{event_id}/category/{category_id}/book [POST]
func (r *rest) CreateBooking(ctx *gin.Context) {
	var bookingParam entity.CreateBookingParam
	if err := ctx.ShouldBindJSON(&bookingParam); err != nil {
		r.httpRespError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	if err := ctx.ShouldBindUri(&bookingParam); err != nil {
		r.httpRespError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	booking, err := r.uc.Booking.Create(ctx.Request.Context(), bookingParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusCreated, "sucessfully create a booking", booking)
}

// @Summary Get Booking
// @Description Get a Booking
// @Security BearerAuth
// @Tags Booking
// @Produce json
// @Param booking_id path string true "booking id param"
// @Success 200 {object} entity.Response{data=[]entity.BookingDetailResponse{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/booking/{booking_id} [GET]
func (r *rest) GetBooking(ctx *gin.Context) {
	var param entity.BookingParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		r.httpRespError(ctx, http.StatusUnprocessableEntity, errors.NewError("request is not valid", err.Error()))
		return
	}

	booking, err := r.uc.Booking.Get(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "sucessfully get booking data", booking)
}
