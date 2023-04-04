package handler

import (
	. "GoServer/Entities"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func checkPassword(password string) error {
	length := len(password)

	if length < 8 {
		return errors.New("password too short")
	}
	if length > 64 {
		return errors.New("password too long")
	}

	var (
		isItHasALowerCase      bool
		isItHasAnUpperCase     bool
		isItHasASpecialCharter bool
		isItHasANumber         bool
	)

	for i := 0; i < length; i++ {
		c := password[i]
		if !isItHasALowerCase && (c >= 'a' && c <= 'z') {
			isItHasALowerCase = true
		} else if !isItHasAnUpperCase && (c >= 'A' && c <= 'Z') {
			isItHasAnUpperCase = true
		} else if !isItHasANumber && (c >= '0' && c <= '9') {
			isItHasANumber = true
		} else if !isItHasASpecialCharter && (c == '&' || c == '$' || c == '@' || c == '!' || c == '-' || c == '_' || c == ' ' || c == '.') {
			isItHasASpecialCharter = true
		} else {
			if !(c >= 'a' && c <= 'z') && !(c >= 'A' && c <= 'Z') && !(c >= '0' && c <= '9') && !(c == '&' || c == '$' || c == '@' || c == '!' || c == '-' || c == '_' || c == ' ' || c == '.') {
				return errors.New("password contains an invalid character")
			}
		}
	}

	if isItHasALowerCase && isItHasAnUpperCase && isItHasASpecialCharter && isItHasANumber {
		return nil
	}

	return errors.New("password has no some characters")
}

func setAllAuthCookies(c *gin.Context, accessToken, longliveToken, email, login string) {
	c.SetCookie("accessToken", accessToken, 60*15, "/", "", false, false)
	c.SetCookie("email", email, 60*60*24*365*100, "/", "", false, true)
	c.SetCookie("isAuth", "true", 60*60*24*365*100, "/", "", false, false)
	c.SetCookie("login", login, 60*60*24*365*100, "/", "", false, true)
	c.SetCookie("longliveToken", longliveToken, 60*60*24*365*100, "/", "", false, true)
}

func clearAllAuthCookies(c *gin.Context) {
	c.SetCookie("accessToken", "", -1, "/", "", false, false)
	c.SetCookie("email", "", -1, "/", "", false, true)
	c.SetCookie("isAuth", "", -1, "/", "", false, false)
	c.SetCookie("login", "", -1, "/", "", false, true)
	c.SetCookie("longliveToken", "", -1, "/", "", false, true)
}

func (handler *Handler) signUp(c *gin.Context) {
	var input UserDTO

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := checkPassword(input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err, Tokens := handler.services.Authorization.CreateUser(c.Request.Context(), input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setAllAuthCookies(c, Tokens.AccessToken, Tokens.LongliveToken, input.Email, input.Login)
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (handler *Handler) signIn(c *gin.Context) {
	var input SignInDTO

	if err := c.BindJSON(&input); err != nil || (input.Email == "" && input.Login == "") {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, email, login, accessToken, longliveToken, err := handler.services.Authorization.SignIn(c.Request.Context(), input)
	if err != nil {
		log.Println(err)
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	setAllAuthCookies(c, accessToken, longliveToken, email, login)

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (handler *Handler) logout(c *gin.Context) {
	clearAllAuthCookies(c)
	c.JSON(http.StatusNoContent, gin.H{})
}
