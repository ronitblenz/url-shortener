package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "url-shortener/handler"
    "url-shortener/model"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

var urlStore *model.URLStore

func setupRouter() *gin.Engine {
    router := gin.Default()
    // Using closures to pass the urlStore to the handlers
    router.POST("/shorten", func(c *gin.Context) {
        handler.ShortenURL(c, urlStore)
    })
    router.GET("/:shortURL", func(c *gin.Context) {
        handler.RedirectURL(c, urlStore)
    })
    router.GET("/metrics", func(c *gin.Context) {
        handler.GetMetrics(c, urlStore)
    })
    return router
}

func TestMain(m *testing.M) {
    urlStore = model.NewURLStore()
    m.Run()
}

func TestShortenURL(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    body := bytes.NewBufferString(`{"url": "http://example.com"}`)
    req, _ := http.NewRequest("POST", "/shorten", body)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "short_url")

    var response map[string]string
    json.Unmarshal(w.Body.Bytes(), &response)
    shortURL := response["short_url"]

    // Ensure the same URL gives the same shortened version
    w = httptest.NewRecorder()
    body = bytes.NewBufferString(`{"url": "http://example.com"}`)
    req, _ = http.NewRequest("POST", "/shorten", body)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "short_url")
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, shortURL, response["short_url"])
}

func TestRedirectURL(t *testing.T) {
    router := setupRouter()

    // Shorten a URL first
    w := httptest.NewRecorder()
    body := bytes.NewBufferString(`{"url": "http://example.com"}`)
    req, _ := http.NewRequest("POST", "/shorten", body)
    router.ServeHTTP(w, req)

    var response map[string]string
    json.Unmarshal(w.Body.Bytes(), &response)
    shortURL := response["short_url"]

    // Now test redirection
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/"+shortURL, nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusMovedPermanently, w.Code)
    assert.Equal(t, "http://example.com", w.Header().Get("Location"))
}

func TestGetMetrics(t *testing.T) {
    router := setupRouter()

    // Shorten multiple URLs
    urls := []string{
        "http://example.com",
        "http://example.com/page1",
        "http://example.org",
        "http://example.net",
        "http://example.net/page",
    }

    for _, url := range urls {
        w := httptest.NewRecorder()
        body := bytes.NewBufferString(`{"url": "` + url + `"}`)
        req, _ := http.NewRequest("POST", "/shorten", body)
        router.ServeHTTP(w, req)
        assert.Equal(t, 200, w.Code)
    }

    // Now get the metrics
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/metrics", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)

    var response map[string]int
    json.Unmarshal(w.Body.Bytes(), &response)

    // Check if the metrics are correct
    assert.Equal(t, 2, response["example.com"])
    assert.Equal(t, 1, response["example.org"])
    assert.Equal(t, 2, response["example.net"])
}
