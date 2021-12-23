package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"price_api/price_server/internal/service"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		startTime := time.Now().Format("2006-01-02 15:04:05")
		startTimeStamp := time.Now().Unix()
		c.Next()

		responseBody := bodyLogWriter.body.String()

		endTime := time.Now().Format("2006-01-02 15:04:05")
		endTimeStamp := time.Now().Unix()

		if c.Request.Method == "POST" {
			c.Request.ParseForm()
		}
		if bodyLogWriter.Status() != http.StatusOK { // not insert log if http status not ok
			return
		}

		accessLogMap := make(map[string]interface{})

		requestUri := c.Request.RequestURI

		accessLogMap["request_time"] = startTime
		accessLogMap["request_uri"] = requestUri
		accessLogMap["request_ua"] = c.Request.UserAgent()
		accessLogMap["request_client_ip"] = c.ClientIP()

		accessLogMap["response_time"] = endTime
		accessLogMap["response"] = responseBody
		accessLogMap["request_timestamp"] = startTimeStamp
		accessLogMap["response_timestamp"] = endTimeStamp

		requestInfoService := service.Svc.RequestInfo()

		if strings.Contains(requestUri, "getPrice") ||
			strings.Contains(requestUri, "getPartyPrice") ||
			strings.Contains(requestUri, "getHistoryPrice") ||
			strings.Contains(requestUri, "getBulkPrices") {
			err := requestInfoService.InsertLogInfo(accessLogMap, 1)
			if err != nil {
				logger.Errorf("insert log info occur err:%v", err)
			}
		} else {
			err := requestInfoService.InsertLogInfo(accessLogMap, 0)
			if err != nil {
				logger.Errorf("insert log info occur err:%v", err)
			}
		}

	}
}
