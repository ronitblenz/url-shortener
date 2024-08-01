package handler

import (
    "net/http"
    "url-shortener/model"
    "github.com/gin-gonic/gin"
)

type URLRequest struct {
    URL string `json:"url" binding:"required"`
}

func ShortenURL(c *gin.Context, urlStore *model.URLStore) {
    var req URLRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    shortURL := urlStore.Save(req.URL)
    c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func RedirectURL(c *gin.Context, urlStore *model.URLStore) {
    shortURL := c.Param("shortURL")
    originalURL, found := urlStore.Get(shortURL)
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
        return
    }
    c.Redirect(http.StatusMovedPermanently, originalURL)
}
