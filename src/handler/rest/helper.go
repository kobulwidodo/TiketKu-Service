package rest

import (
	"fmt"
	"go-clean/src/business/entity"
	"go-clean/src/lib/appcontext"
	"go-clean/src/lib/errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func (r *rest) httpRespSuccess(ctx *gin.Context, code int, message string, data interface{}) {
	c := ctx.Request.Context()
	resp := entity.Response{
		Meta: entity.Meta{
			Message:   message,
			Code:      code,
			IsError:   false,
			RequestID: appcontext.GetRequestID(c),
		},
		Data: data,
	}
	ctx.JSON(code, resp)
}

func (r *rest) httpRespError(ctx *gin.Context, code int, err error) {
	c := ctx.Request.Context()
	r.log.Error(c, err)
	resp := entity.Response{
		Meta: entity.Meta{
			Message:   err.Error(),
			Code:      code,
			IsError:   true,
			RequestID: appcontext.GetRequestID(c),
		},
		Data: nil,
	}
	ctx.AbortWithStatusJSON(code, resp)
}

func (r *rest) VerifyUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.NewError("authentication failed", "empty token"))
		return
	}

	var tokenString string
	_, err := fmt.Sscanf(authHeader, "Bearer %v", &tokenString)
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.NewError("authentication failed", "invalid token"))
		return
	}

	token, err := r.ValidateToken(tokenString)
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, err)
		return
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.NewError("authentication failed", "failed to claim token"))
		return
	}

	user := entity.User{}
	user, err = r.uc.User.GetById(uint(claim["id"].(float64)))
	if err != nil {
		r.httpRespError(ctx, http.StatusUnauthorized, errors.NewError("authentication failed", "error while getting user"))
		return
	}

	c := ctx.Request.Context()
	c = r.auth.SetUserAuthInfo(c, user.ConvertToAuthUser(), tokenString)
	ctx.Request = ctx.Request.WithContext(c)

	ctx.Next()
}

func (r *rest) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.NewError("authentication failed", "token invalid")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (r *rest) addFieldsToContext(ctx *gin.Context) {
	reqid := ctx.GetHeader(entity.XRequestId)
	if reqid == "" {
		reqid = uuid.New().String()
	}

	c := ctx.Request.Context()
	c = appcontext.SetRequestID(c, reqid)
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()
}
