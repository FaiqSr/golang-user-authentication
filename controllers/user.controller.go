package controllers

import (
	"fmt"
	"golang-user-authentication/dto"
	"golang-user-authentication/helpers"
	"golang-user-authentication/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {

	var user dto.UserCreateUserRequest

	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	err = user.ValidateUserCreateRequest()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res, err := models.CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, res)
}

func Login(c echo.Context) error {
	var loginRequest dto.UserLoginRequest
	err := c.Bind(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	err = loginRequest.ValidateUserLoginRequest()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res, err := models.LoginUser(&loginRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, res)

}

func GetUser(c echo.Context) error {
	userId := helpers.GetUserIdFromJwt(c)

	res, err := models.GetUserInformation(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func UpdateUser(c echo.Context) error {

	var user dto.UserUpdateRequest

	userId := helpers.GetUserIdFromJwt(c)

	err := c.Bind(&user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	err = user.ValidateUserUpdateRequest()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res, err := models.UpdateUser(&user, userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func LogoutUser(c echo.Context) error {
	userId := helpers.GetUserIdFromJwt(c)

	res, err := models.LogoutUser(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, res)

}
