package main

import (
    "url-shortener/handler"
    "url-shortener/model"
    "github.com/gin-gonic/gin"
    "os"
)

func main() {
    redisHost := os.Getenv("REDIS_HOST")
    if redisHost == "" {
        redisHost = "localhost:6379"
    }

    urlStore := model.NewURLStore(redisHost)
    router := gin.Default()

    router.POST("/shorten", func(c *gin.Context) {
        handler.ShortenURL(c, urlStore)
    })
    router.GET("/:shortURL", func(c *gin.Context) {
        handler.RedirectURL(c, urlStore)
    })
    router.GET("/metrics", func(c *gin.Context) {
        handler.GetMetrics(c, urlStore)
    })

    router.Run(":8080")
}
