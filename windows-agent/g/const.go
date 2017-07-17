package g

import (
	"time"
)

// changelog:
// 1.0.1
// 1.0.2 add process support
// 51idc-1.0.3 support netstat metrics
// add sysinfo err catch
// 51idc-1.0.5 support design process monitor & design port monitor
const (
	VERSION          = "51idc-1.0.5"
	COLLECT_INTERVAL = time.Second
)
