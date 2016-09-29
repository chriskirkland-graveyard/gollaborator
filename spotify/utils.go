package spotify

import "time"

var GlobalRateLimiter = time.Tick(100 * time.Millisecond)
