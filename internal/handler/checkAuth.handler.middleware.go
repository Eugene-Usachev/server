package handler

import (
	"GoServer/pkg/fasthttp_utils"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/gofiber/fiber/v2"
)

// TODO r
/*func (handler *handler) refreshTokens(c *fiber.Ctx) bool {
	longliveToken := c.Cookies("longliveToken")
	if longliveToken == "" {
		NewErrorResponse(c, fiber.StatusUnauthorized, "empty auth token")
		return false
	}

	password, err := jwt.ParseLongliveToken(longliveToken)
	if err != nil {
		NewErrorResponse(c, fiber.StatusUnauthorized, "invalid auth token")
		return false
	}

	email := c.Cookies("email")
	if email == "" {
		NewErrorResponse(c, fiber.StatusUnauthorized, "empty email cookie")
		return false
	}

	var (
		id          uint
		accessToken string
	)

	id, accessToken, longliveToken, err = handler.services.RefreshToken(c.Context(), password, email)
	if err != nil {
		NewErrorResponse(c, fiber.StatusUnauthorized, "invalid auth token")
		return false
	}

	login := c.Cookies("login")
	setAllAuthCookies(c, accessToken, longliveToken, email, login)

	c.Locals("userId", id)
	return true
}*/

func (handler *Handler) CheckAuth(c *fiber.Ctx) error {
	accessToken := fastbytes.B2S(fasthttp_utils.GetAuthorizationHeader(c.Context()))
	if accessToken == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	} else {
		userId, err := handler.accessConverter.ParseToken(accessToken)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		} else {
			id := fastbytes.B2U(userId)
			c.Locals("userId", id)
		}
	}
	return c.Next()
}

// TODO r
// SetTokensInFirst is using to set tokens when user get a page. This is important, because, in first,
// we need to refreshTokens even on GET request and, in second, to use websocket correct.
/*func (handler *handler) SetTokensInFirst(ctx *fiber.Ctx) error {
	accessToken := ctx.Cookies("accessToken")
	if accessToken == "" {
		handler.refreshInFirst(ctx)
	} else {
		userId, err := jwt.ParseAccessToken(accessToken)
		if err != nil {
			handler.refreshInFirst(ctx)
			ctx.Locals("userId", userId)
		}
	}
	return ctx.Next()
}

func (handler *handler) refreshInFirst(ctx *fiber.Ctx) bool {
	longliveToken := ctx.Cookies("longliveToken")
	if longliveToken == "" {
		return false
	}
	password, err := jwt.ParseLongliveToken(longliveToken)
	if err != nil {
		return false
	}

	email := ctx.Cookies("email")
	if email == "" {
		return false
	}

	id, accessToken, longliveToken, err := handler.services.RefreshToken(ctx.Context(), password, email)
	if err != nil {
		return false
	}
	login := ctx.Cookies("login")
	setAllAuthCookies(ctx, accessToken, longliveToken, email, login)

	ctx.Locals("userId", id)
	return true
}
*/
