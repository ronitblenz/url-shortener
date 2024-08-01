package handler

import (
    "net/http"
    "url-shortener/model"
    "github.com/gin-gonic/gin"
)

var urlStore = model.NewURLStore()

type URLRequest struct {
    URL string `json:"url" binding:"required"`
}

func ShortenURL(c *gin.Context) {
    var req URLRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    shortURL := urlStore.Save(req.URL)
    c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func RedirectURL(c *gin.Context) {
    shortURL := c.Param("shortURL")
    originalURL, found := urlStore.Get(shortURL)
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }
    c.Redirect(http.StatusMovedPermanently, originalURL)
}
