package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"xray-geo-exporter/config"
	"xray-geo-exporter/geo"
)

type Payload struct {
	GeoData   *geo.GeoData `json:"geo_data"`
	Hostname  string       `json:"hostname"`
	Timestamp time.Time    `json:"timestamp"`
}

type Sender struct {
	client     *http.Client
	webhookURL string
}

func NewSender(webhookURL string, timeout time.Duration) *Sender {
	return &Sender{
		client: &http.Client{
			Timeout: timeout,
		},
		webhookURL: webhookURL,
	}
}

func (s *Sender) Send(geoData *geo.GeoData) error {
	payload := Payload{
		GeoData:   geoData,
		Hostname:  config.Hostname,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status code: %d", resp.StatusCode)
	}

	return nil
}
