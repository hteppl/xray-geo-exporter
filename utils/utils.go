package utils

import (
	"log"
	"sync"
	"time"
	"xray-geo-exporter/config"
	"xray-geo-exporter/geo"
	"xray-geo-exporter/ratelimit"
	"xray-geo-exporter/storage"
	"xray-geo-exporter/webhook"

	"github.com/hpcloud/tail"
)

type LogMonitor struct {
	ipStorage     *storage.IPStorage
	rateLimiter   *ratelimit.RateLimiter
	geoFetcher    *geo.Fetcher
	webhookSender *webhook.Sender
	ipQueue       chan string
	wg            sync.WaitGroup
}

func NewLogMonitor() *LogMonitor {
	return &LogMonitor{
		ipStorage:     storage.NewIPStorage(time.Duration(config.IPTTLMinutes) * time.Minute),
		rateLimiter:   ratelimit.NewRateLimiter(config.RateLimitPerMinute),
		geoFetcher:    geo.NewFetcher(config.HTTPTimeout),
		webhookSender: webhook.NewSender(config.WebhookURL, config.WebhookTimeout),
		ipQueue:       make(chan string, config.QueueSize),
	}
}

func (m *LogMonitor) Start() {
	// Start IP processor
	m.wg.Add(1)
	go m.processIPQueue()

	// Start log monitor
	t, err := tail.TailFile(config.LogFile, tail.Config{
		Follow:    true,
		ReOpen:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
	})
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	for line := range t.Lines {
		m.handleLogEntry(line.Text)
	}
}

func (m *LogMonitor) handleLogEntry(line string) {
	ip := config.IpRegex.FindString(line)
	if ip == "" {
		return
	}

	// Check if IP is already processed (within TTL)
	if !m.ipStorage.Add(ip) {
		return
	}

	// Add to processing queue
	select {
	case m.ipQueue <- ip:
	default:
		log.Printf("IP queue full, dropping IP: %s", ip)
	}
}

func (m *LogMonitor) processIPQueue() {
	defer m.wg.Done()

	for ip := range m.ipQueue {
		// Rate limit before making request
		m.rateLimiter.Wait()

		// Fetch geo data
		geoData, err := m.geoFetcher.FetchGeoData(ip)
		if err != nil {
			log.Printf("Failed to fetch geo data for IP %s: %v", ip, err)
			continue
		}

		// Send to webhook
		if err := m.webhookSender.Send(geoData); err != nil {
			log.Printf("Failed to send webhook for IP %s: %v", ip, err)
			continue
		}

		log.Printf("Successfully processed IP %s from %s, %s", ip, geoData.City, geoData.Country)
	}
}

func (m *LogMonitor) Stop() {
	close(m.ipQueue)
	m.wg.Wait()
}

func StartLogMonitor() {
	monitor := NewLogMonitor()
	monitor.Start()
}
