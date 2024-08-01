package model

import (
    "net/url"
    "sort"
    "sync"
)

const (
    alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    base     = int64(len(alphabet))
)

type URLStore struct {
    sync.RWMutex
    urls     map[string]string
    reverse  map[string]string // for reverse lookup
    counter  int64
}

func NewURLStore() *URLStore {
    return &URLStore{
        urls:    make(map[string]string),
        reverse: make(map[string]string),
        counter: 1,
    }
}

func (s *URLStore) Save(originalURL string) string {
    s.Lock()
    defer s.Unlock()

    // Check if the URL already exists
    if shortURL, exists := s.reverse[originalURL]; exists {
        return shortURL
    }

    // Increment counter and encode it
    shortURL := encode(s.counter)
    s.counter++
    s.urls[shortURL] = originalURL
    s.reverse[originalURL] = shortURL
    return shortURL
}

func (s *URLStore) Get(shortURL string) (string, bool) {
    s.RLock()
    defer s.RUnlock()
    originalURL, found := s.urls[shortURL]
    return originalURL, found
}

func (s *URLStore) GetTopDomains(limit int) map[string]int {
    s.RLock()
    defer s.RUnlock()

    domainCount := make(map[string]int)

    for _, originalURL := range s.urls {
        parsedURL, err := url.Parse(originalURL)
        if err == nil {
            domain := parsedURL.Host
            domainCount[domain]++
        }
    }

    // Create a slice of domain names and sort by frequency
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

    // Prepare the final result with a limit
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
