package domain

import (
	"testing"
	"time"
)

func TestConcert_IsNew(t *testing.T) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	window := 7 * 24 * time.Hour

	tests := []struct {
		name         string
		discoveredAt time.Time
		want         bool
	}{
		{"just now", now, true},
		{"exactly at window boundary", now.Add(-window), true},
		{"just past window boundary", now.Add(-window - time.Second), false},
		{"long ago", now.AddDate(0, -1, 0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Concert{DiscoveredAt: tt.discoveredAt}
			if got := c.IsNew(now, window); got != tt.want {
				t.Errorf("IsNew(discoveredAt=%v) = %v, want %v", tt.discoveredAt, got, tt.want)
			}
		})
	}
}
