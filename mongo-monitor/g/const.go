package g

import (
	"time"
)

// changelog:
// 0.0.1: first version
// 0.0.2: fix auth bug
// 0.0.3: fix pannic in 2.6
const (
	VERSION          = "0.0.3"
	COLLECT_INTERVAL = time.Second
)
