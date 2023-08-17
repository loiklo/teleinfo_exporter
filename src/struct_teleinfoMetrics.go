package main

import "sync"

type TeleinfoMetrics struct {
	mu     sync.Mutex
	metric map[string]string
}

// Safe setter
func (tm *TeleinfoMetrics) Set(label string, data string) {
	tm.mu.Lock()
	tm.metric[label] = data
	tm.mu.Unlock()
}

// Safe getter
func (tm *TeleinfoMetrics) Get(label string) string {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.metric[label]
}
