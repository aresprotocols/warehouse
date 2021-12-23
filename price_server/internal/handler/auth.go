package handler

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/constant"
	"price_api/price_server/internal/service"
	"price_api/price_server/internal/vo"
)

func HandleAuth(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	var user vo.AdminUser
	err := context.ShouldBind(&user)
	if err != nil {
		logger.WithError(err).Errorf("bind user occur error")
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(user.User) == 0 {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	if len(user.Password) == 0 {
		response.Code = constant.PARAM_NOT_TRUE_ERROR
		response.Message = constant.MSG_PARAM_NOT_TRUE
		context.JSON(http.StatusBadRequest, response)
		return
	}

	authService := service.Svc.Auth()

	if !authService.ValidateUserAndPassword(user) {
		response.Code = constant.CHECK_USER_ERROR
		response.Message = constant.MSG_CHECK_USER_ERROR
		context.JSON(http.StatusUnauthorized, response)
		return
	}
	authToken, err := authService.GenerateToken(user)
	if err != nil {
		logger.WithError(err).Error("generate jwt token occur error")
		response.Code = constant.ERROR
		response.Message = err.Error()
		context.JSON(http.StatusInternalServerError, response)
	}
	response.Data = authToken
	context.JSON(http.StatusOK, response)

}
