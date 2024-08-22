package rest

import (
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get Seat
// @Description Get a Seat
// @Security BearerAuth
// @Tags Seat
// @Produce json
// @Param event_id path int true "event id param"
// @Param category_id path int true "category id param"
// @Success 200 {object} entity.Response{data=[]entity.Event{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/event/{event_id}/category/{category_id}/seat [GET]
func (r *rest) GetListSeat(ctx *gin.Context) {
	var seatParam entity.SeatParam
	if err := ctx.ShouldBindUri(&seatParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, errors.NewError("request cannot be processed", err.Error()))
		return
	}

	event, err := r.uc.Seat.GetList(seatParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get seat list", event)
}
