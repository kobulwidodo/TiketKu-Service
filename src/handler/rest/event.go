package rest

import (
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get List Event
// @Description Get List All Event
// @Security BearerAuth
// @Tags Event
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.Event{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/event [GET]
func (r *rest) GetListEvent(ctx *gin.Context) {
	events, err := r.uc.Event.GetList(entity.EventParam{})
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get list all event", events)
}

// @Summary Get Event
// @Description Get a Event
// @Security BearerAuth
// @Tags Event
// @Produce json
// @Param event_id path int true "event id param"
// @Success 200 {object} entity.Response{data=[]entity.Event{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/event/{event_id} [GET]
func (r *rest) GetEvent(ctx *gin.Context) {
	var eventParam entity.EventParam
	if err := ctx.ShouldBindUri(&eventParam); err != nil {
		r.httpRespError(ctx, http.StatusBadRequest, errors.NewError("request cannot be processed", err.Error()))
		return
	}

	event, err := r.uc.Event.Get(eventParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get a event", event)
}
