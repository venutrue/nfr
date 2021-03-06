package client

import (
	"context"
	"encoding/json"
	"net/url"
)

// AlertsResponse represents response for /events call.
type AlertsResponse struct {
	Follow string `json:"follow"`
	More   bool   `json:"more"`
	After  string `json:"after,omitempty"`
	Before string `json:"before,omitempty"`

	Alerts  []Alert           `json:"alerts,omitempty"`
	Threats map[string]Threat `json:"threats,omitempty"`
}

// Alert provides result of AlphaSOC Engine analysis, which was found to be threat.
type Alert struct {
	EventType string          `json:"eventType"`
	Event     json.RawMessage `json:"event"`
	IPEvent   IPEntry         `json:"-"`
	DNSEvent  DNSEntry        `json:"-"`
	Threats   []string        `json:"threats"`
	Wisdom    struct {
		Flags []string `json:"flags"`
	} `json:"wisdom"`
}

// Threat provides more details about threat,
// like human-readable description.
type Threat struct {
	Title    string `json:"title"`
	Severity int    `json:"severity"`
	Policy   bool   `json:"policy"`
}

// Alerts returns AlphaSOC events that informs about potential risk.
func (c *AlphaSOCClient) Alerts(follow string) (*AlertsResponse, error) {
	if c.key == "" {
		return nil, ErrNoAPIKey
	}
	query := url.Values{}
	if follow != "" {
		query.Add("follow", follow)
	}
	resp, err := c.get(context.Background(), "alerts", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r AlertsResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	for i := range r.Alerts {
		switch r.Alerts[i].EventType {
		case "dns":
			if err := json.Unmarshal(r.Alerts[i].Event, &r.Alerts[i].DNSEvent); err != nil {
				return nil, err
			}
		case "ip":
			if err := json.Unmarshal(r.Alerts[i].Event, &r.Alerts[i].IPEvent); err != nil {
				return nil, err
			}
		}
	}

	return &r, nil
}
