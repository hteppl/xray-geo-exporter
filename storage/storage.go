package storage

import (
	"sync"
	"time"
)

type IPEntry struct {
	IP        string
	Timestamp time.Time
}

type IPStorage struct {
	mu      sync.RWMutex
	entries map[string]time.Time
	ttl     time.Duration
}

func NewIPStorage(ttl time.Duration) *IPStorage {
	storage := &IPStorage{
		entries: make(map[string]time.Time),
		ttl:     ttl,
	}

	go storage.cleanupExpired()

	return storage
}

func (s *IPStorage) Add(ip string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if expiry, exists := s.entries[ip]; exists {
		if time.Now().Before(expiry) {
			return false
		}
	}

	s.entries[ip] = time.Now().Add(s.ttl)
	return true
}

func (s *IPStorage) Exists(ip string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if expiry, exists := s.entries[ip]; exists {
		return time.Now().Before(expiry)
	}
	return false
}

func (s *IPStorage) cleanupExpired() {
	ticker := time.NewTicker(s.ttl / 10) // Cleanup interval is 1/10th of TTL
	if s.ttl < 10*time.Minute {
		ticker = time.NewTicker(1 * time.Minute) // Minimum 1 minute
	}
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for ip, expiry := range s.entries {
			if now.After(expiry) {
				delete(s.entries, ip)
			}
		}
		s.mu.Unlock()
	}
}
