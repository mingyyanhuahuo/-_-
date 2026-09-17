package middleware

import (
	"lostfound/pkg/redisdb"

	"github.com/gin-gonic/gin"
)

func FrequentWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		err, ErrCode := redisdb.RequestsTimeLimit(ip)
		if err != nil {

		}
		if ErrCode != nil {

		}
	}
}
