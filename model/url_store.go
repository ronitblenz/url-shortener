package model

import (
    "context"
    "github.com/go-redis/redis/v8"
    "net/url"
    "sort"
    "strconv"
    "sync"
)

const (
    alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    base     = int64(len(alphabet))
)

type URLStore struct {
    sync.RWMutex
    client  *redis.Client
    counter int64
}

func NewURLStore(redisAddr string) *URLStore {
    client := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    return &URLStore{
        client:  client,
        counter: 1,
    }
}

func (s *URLStore) Save(originalURL string) (string, error) {
    s.Lock()
    defer s.Unlock()
    ctx := context.Background()

    // Check if the URL already exists
    shortURL, err := s.client.Get(ctx, originalURL).Result()
    if err == nil {
        return shortURL, nil
    }

    // Increment counter and encode it
    shortURL = encode(s.counter)
    s.counter++

    // Store short URL and original URL
    err = s.client.Set(ctx, shortURL, originalURL, 0).Err()
    if err != nil {
        return "", err
    }
    err = s.client.Set(ctx, originalURL, shortURL, 0).Err()
    if err != nil {
        return "", err
    }

    // Update counter in Redis
    err = s.client.Set(ctx, "url-counter", strconv.FormatInt(s.counter, 10), 0).Err()
    if err != nil {
        return "", err
    }

    return shortURL, nil
}

func (s *URLStore) Get(shortURL string) (string, bool) {
    s.RLock()
    defer s.RUnlock()
    ctx := context.Background()
    originalURL, err := s.client.Get(ctx, shortURL).Result()
    if err != nil {
        return "", false
    }
    return originalURL, true
}

func (s *URLStore) GetTopDomains(limit int) map[string]int {
    s.RLock()
    defer s.RUnlock()
    ctx := context.Background()
    domainCount := make(map[string]int)

    keys, err := s.client.Keys(ctx, "*").Result()
    if err != nil {
        return domainCount
    }

    for _, key := range keys {
        originalURL, err := s.client.Get(ctx, key).Result()
        if err == nil {
            domain := extractDomain(originalURL)
            if domain != "" {
                domainCount[domain]++
            }
        }
    }

    // Convert map to a slice and sort it by frequency
    type kv struct {
        Key   string
        Value int
    }

    var sortedDomains []kv
    for k, v := range domainCount {
        sortedDomains = append(sortedDomains, kv{k, v})
    }

    sort.Slice(sortedDomains, func(i, j int) bool {
        return sortedDomains[i].Value > sortedDomains[j].Value
    })

    // Limit the results
    topDomains := make(map[string]int)
    for i, domain := range sortedDomains {
        if i >= limit {
            break
        }
        topDomains[domain.Key] = domain.Value
    }

    return topDomains
}

func encode(num int64) string {
    if num == 0 {
        return string(alphabet[0])
    }

    encoded := ""
    for num > 0 {
        remainder := num % base
        encoded = string(alphabet[remainder]) + encoded
        num = num / base
    }
    return encoded
}

func extractDomain(urlStr string) string {
    u, err := url.Parse(urlStr)
    if err != nil || u.Host == "" {
        return ""
    }
    return u.Host
}
