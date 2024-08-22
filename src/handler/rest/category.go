package rest

import (
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Get List Category
// @Description Get List All Category
// @Security BearerAuth
// @Tags Category
// @Produce json
// @Param event_id path int true "event id param"
// @Success 200 {object} entity.Response{data=[]entity.Category{}}
// @Failure 400 {object} entity.Response{}
// @Failure 401 {object} entity.Response{}
// @Failure 404 {object} entity.Response{}
// @Failure 500 {object} entity.Response{}
// @Router /api/v1/event/{event_id}/category [GET]
func (r *rest) GetListCategory(ctx *gin.Context) {
	var categoryParam entity.CategoryParam
	if err := ctx.ShouldBindUri(&categoryParam); err != nil {
		r.httpRespError(ctx, http.StatusUnprocessableEntity, errors.NewError("request cannot be processed", err.Error()))
		return
	}

	categories, err := r.uc.Category.GetList(categoryParam)
	if err != nil {
		r.httpRespError(ctx, http.StatusInternalServerError, err)
		return
	}

	r.httpRespSuccess(ctx, http.StatusOK, "successfully get list all categories", categories)
}
