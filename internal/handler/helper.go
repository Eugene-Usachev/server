package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func NewErrorResponse(ctx *fiber.Ctx, statusCode int, message string) error {
	return ctx.Status(statusCode).JSON(fiber.Map{
		"message": message,
	})
}

func setCookie(c *fiber.Ctx, key, value string, maxAge int) {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.Value = value
	cookie.MaxAge = maxAge
	cookie.Path = "/"
	cookie.HTTPOnly = true
	c.Cookie(cookie)
}

func checkPassword(password string) error {
	length := len(password)

	if length < 8 {
		return errors.New("password too short")
	}
	if length > 64 {
		return errors.New("password too long")
	}

	//var (
	//	isItHasALowerCase      bool
	//	isItHasAnUpperCase     bool
	//	isItHasASpecialCharter bool
	//	isItHasANumber         bool
	//)
	//
	//for i := 0; i < length; i++ {
	//	c := password[i]
	//	if !isItHasALowerCase && (c >= 'a' && c <= 'z') {
	//		isItHasALowerCase = true
	//	} else if !isItHasAnUpperCase && (c >= 'A' && c <= 'Z') {
	//		isItHasAnUpperCase = true
	//	} else if !isItHasANumber && (c >= '0' && c <= '9') {
	//		isItHasANumber = true
	//	} else if !isItHasASpecialCharter && (c == '&' || c == '$' || c == '@' || c == '!' || c == '-' || c == '_' || c == ' ' || c == '.') {
	//		isItHasASpecialCharter = true
	//	} else {
	//		if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && !(c == '&' || c == '$' || c == '@' || c == '!' || c == '-' || c == '_' || c == ' ' || c == '.') {
	//			return errors.New("password contains an invalid character")
	//		}
	//	}
	//}
	//
	//if isItHasALowerCase && isItHasAnUpperCase && isItHasASpecialCharter && isItHasANumber {
	//	return nil
	//}
	//
	//return errors.New("password has no some characters")
	return nil
}
