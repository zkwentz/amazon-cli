package amazon

import (
	"testing"
)

func TestGetRandomUserAgent(t *testing.T) {
	// Test that getRandomUserAgent returns a non-empty string
	userAgent := getRandomUserAgent()
	if userAgent == "" {
		t.Error("getRandomUserAgent returned empty string")
	}
}

func TestGetRandomUserAgentReturnsValidUserAgent(t *testing.T) {
	// Test that the returned user agent is one from our list
	userAgent := getRandomUserAgent()

	found := false
	for _, ua := range userAgents {
		if ua == userAgent {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("getRandomUserAgent returned unexpected user agent: %s", userAgent)
	}
}

func TestGetRandomUserAgentDistribution(t *testing.T) {
	// Test that calling getRandomUserAgent multiple times returns different values
	// This tests the randomness aspect
	calls := 100
	results := make(map[string]int)

	for i := 0; i < calls; i++ {
		ua := getRandomUserAgent()
		results[ua]++
	}

	// With 100 calls and 10 user agents, we should see at least 2 different user agents
	// (this is a probabilistic test, but the chances of getting only 1 UA in 100 calls is extremely low)
	if len(results) < 2 {
		t.Errorf("getRandomUserAgent showed poor randomness: only %d unique user agents in %d calls", len(results), calls)
	}
}

func TestUserAgentsSliceHasCorrectLength(t *testing.T) {
	// Test that we have exactly 10 user agents as specified
	expectedCount := 10
	if len(userAgents) != expectedCount {
		t.Errorf("userAgents slice has %d entries, expected %d", len(userAgents), expectedCount)
	}
}

func TestUserAgentsSliceContainsValidStrings(t *testing.T) {
	// Test that all user agents are non-empty and appear valid
	for i, ua := range userAgents {
		if ua == "" {
			t.Errorf("userAgents[%d] is empty", i)
		}

		// User agents should be reasonably long (at least 50 characters)
		if len(ua) < 50 {
			t.Errorf("userAgents[%d] appears to be invalid (too short): %s", i, ua)
		}

		// All modern browser user agents contain "Mozilla"
		if len(ua) > 0 && len(ua) >= 7 {
			found := false
			for j := 0; j <= len(ua)-7; j++ {
				if ua[j:j+7] == "Mozilla" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("userAgents[%d] does not contain 'Mozilla': %s", i, ua)
			}
		}
	}
}
