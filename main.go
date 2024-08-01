package main

import (
    "url-shortener/handler"
    "url-shortener/model"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    urlStore := model.NewURLStore()

    r.POST("/shorten", func(c *gin.Context) {
        handler.ShortenURL(c, urlStore)
    })
    r.GET("/:shortURL", func(c *gin.Context) {
        handler.RedirectURL(c, urlStore)
    })
    r.GET("/metrics", func(c *gin.Context) {
        handler.GetMetrics(c, urlStore)
    })

    r.Run(":8080")
}
