package api

import (
	"chym/stream/backend/protocols"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(c *gin.Context) {
	proxyReq := protocols.ProxyReq{}
	err := c.ShouldBindJSON(&proxyReq)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, GenResponse(nil, PARAMETER_ERROR, "FAILED"))
		return
	}

}
