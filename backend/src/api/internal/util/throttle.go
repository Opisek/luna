package util

import "time"

type Throttle struct {
	failedAttempts    map[string]int
	lastFailedAttempt map[string]time.Time
}

func NewRequestThrottle() *Throttle {
	return &Throttle{
		failedAttempts:    make(map[string]int),
		lastFailedAttempt: make(map[string]time.Time),
	}
}

func (t *Throttle) RecordFailedAttempt(ip string) {
	if _, exists := t.failedAttempts[ip]; !exists {
		t.failedAttempts[ip] = 1
	} else {
		t.failedAttempts[ip]++
	}
	t.lastFailedAttempt[ip] = time.Now()
}

func (t *Throttle) RecordSuccessfulAttempt(ip string) {
	delete(t.failedAttempts, ip)
	delete(t.lastFailedAttempt, ip)
}

func (t *Throttle) GetFailedAttempts(ip string) int {
	lastFail, ok := t.lastFailedAttempt[ip]
	if ok && time.Since(lastFail) > 1*time.Minute {
		delete(t.failedAttempts, ip)
		delete(t.lastFailedAttempt, ip)
	} else if ok {
		return t.failedAttempts[ip]
	}
	return 0
}

func (t *Throttle) CleanStaleEntries() {
	for ip, lastFail := range t.lastFailedAttempt {
		if time.Since(lastFail) > 1*time.Minute {
			delete(t.failedAttempts, ip)
			delete(t.lastFailedAttempt, ip)
		}
	}
}
