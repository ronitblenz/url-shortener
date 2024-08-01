# Go URL Shortener

A simple URL shortener service built with Go and Gin. This service allows users to shorten URLs, redirect to the original URL, and track domain metrics.

```This is made for the assignment submission for InfraCloud Technologies```

## Getting Started

### Docker Setup

```bash
docker-compose up --build
```

### Local Setup

1. **Clone the repository:**

```bash
   git clone https://github.com/ronitblenz/url-shortener.git
   cd url-shortener
```

2. **Install dependencies:**

```bash
   go mod tidy
```

### Running the Application

To run the URL shortener service, execute:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

### API Endpoints

- **Shorten a URL:**

  `POST /shorten`

  Request Body:
  ```json
  {
    "url": "http://example.com"
  }
  ```

  Response:
  ```json
  {
    "short_url": "shortenedURL"
  }
  ```

For example: We are generating "1" for the first URL, "2" for the next, and so on.

- **Redirect to Original URL:**

  `GET /shortenedURL`

  Redirects to the original URL.

- **Get Metrics:**

  `GET /metrics`

  Returns a JSON object containing the counts of shortened URLs grouped by domain.

### Running Tests

To run the tests, use the following command:

```bash
go test ./tests
```

This will execute all the test cases in the `tests` package and display the results.

### Docker Image Link : 

```bash
https://hub.docker.com/layers/ronitblenz/url-shortener/latest/images/latest
```
