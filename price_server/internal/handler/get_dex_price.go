package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"price_api/price_server/internal/vo"
)

func HandleGetDexPrice(context *gin.Context) {
	response := vo.RESPONSE{Code: 0, Message: "OK"}

	aresShowInfo := handle.fetcher.GetDexPrice()

	response.Data = aresShowInfo
	context.JSON(http.StatusOK, response)
}
