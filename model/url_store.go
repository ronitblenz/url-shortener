package model

import (
    "crypto/sha1"
    "encoding/hex"
    "sync"
)

type URLStore struct {
    store map[string]string
    lock  sync.RWMutex
}

func NewURLStore() *URLStore {
    return &URLStore{store: make(map[string]string)}
}

func (s *URLStore) Save(url string) string {
    s.lock.Lock()
    defer s.lock.Unlock()
    hash := sha1.New()
    hash.Write([]byte(url))
    shortURL := hex.EncodeToString(hash.Sum(nil))[:8]
    s.store[shortURL] = url
    return shortURL
}

func (s *URLStore) Get(shortURL string) (string, bool) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    url, found := s.store[shortURL]
    return url, found
}

func (s *URLStore) GetTopDomains(limit int) map[string]int {
    domainCount := make(map[string]int)
    for _, url := range s.store {
        domain := extractDomain(url)
        domainCount[domain]++
    }

    topDomains := make(map[string]int)
    // Extract top N domains based on count
    for domain, count := range domainCount {
        topDomains[domain] = count
        if len(topDomains) == limit {
            break
        }
    }
    return topDomains
}

func extractDomain(url string) string {
    // Logic to extract domain from the URL
    return url // Placeholder
}
