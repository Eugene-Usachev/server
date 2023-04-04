package handler

import (
	"GoServer/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (handler *Handler) refresh(ctx *gin.Context) bool {
	longliveToken, err := ctx.Cookie("longliveToken")
	if err != nil || longliveToken == "" {
		NewErrorResponse(ctx, http.StatusUnauthorized, "empty auth token")
		return false
	}
	password, err := jwt.ParseLongliveToken(longliveToken)
	if err != nil {
		NewErrorResponse(ctx, http.StatusUnauthorized, "empty auth token")
		return false
	}

	email, err := ctx.Cookie("email")
	if err != nil || email == "" {
		NewErrorResponse(ctx, http.StatusUnauthorized, "empty email cookie")
		return false
	}

	id, accessToken, longliveToken, err := handler.services.RefreshToken(ctx.Request.Context(), password, email)
	if err != nil {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth token")
		return false
	}
	login, _ := ctx.Cookie("login")
	setAllAuthCookies(ctx, accessToken, longliveToken, email, login)

	ctx.Set("userId", id)
	return true
}

func (handler *Handler) CheckAuth(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("accessToken")
	if err != nil || accessToken == "" {
		if !handler.refresh(ctx) {
			return
		}
	} else {
		var userId uint
		userId, err = jwt.ParseAccessToken(accessToken)
		if err != nil {
			if !handler.refresh(ctx) {
				return
			}
		} else {
			ctx.Set("userId", userId)
		}
	}
}

// SetTokensInFirst is using to set tokens when user get a page. This is important, because, in first,
// we need to refresh token even on GET request and, in second, to use websocket correct.
func (handler *Handler) SetTokensInFirst(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("accessToken")
	if err != nil || accessToken == "" {
		handler.refreshInFirst(ctx)
	} else {
		userId, err := jwt.ParseAccessToken(accessToken)
		if err != nil {
			handler.refreshInFirst(ctx)
			ctx.Set("userId", userId)
		}
	}
}

func (handler *Handler) refreshInFirst(ctx *gin.Context) bool {
	longliveToken, err := ctx.Cookie("longliveToken")
	if err != nil || longliveToken == "" {
		return false
	}
	password, err := jwt.ParseLongliveToken(longliveToken)
	if err != nil {
		return false
	}

	email, err := ctx.Cookie("email")
	if err != nil || email == "" {
		return false
	}

	id, accessToken, longliveToken, err := handler.services.RefreshToken(ctx.Request.Context(), password, email)
	if err != nil {
		return false
	}
	login, _ := ctx.Cookie("login")
	setAllAuthCookies(ctx, accessToken, longliveToken, email, login)

	ctx.Set("userId", id)
	return true
}
