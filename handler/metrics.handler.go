package handler

import (
    "net/http"
    "url-shortener/model"
    "github.com/gin-gonic/gin"
)

func GetMetrics(c *gin.Context) {
    metrics := urlStore.GetTopDomains(3)
    c.JSON(http.StatusOK, metrics)
}
