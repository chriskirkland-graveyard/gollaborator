package spotify

import "time"

var GlobalRateLimiter = time.Tick(30 * time.Millisecond)
