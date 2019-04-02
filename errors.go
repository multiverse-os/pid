package pid

import (
	"fmt"
)

var (
	// Segregated to simplify localization
	errStalePid    = fmt.Errorf("stale pid check open failed:", err)
	errFileLocked  = fmt.Errorf("pid file locked by running process:", err)
	errOpenFailed  = fmt.Errorf("failed to open pid file:", err)
	errWriteFailed = fmt.Errorf("failed to write pid to file:", err)
	errShortWrite  = fmt.Errorf("short write failure writing pid:", err)
	errCleanFailed = fmt.Errorf("failed to clean up pid file:", err)
)
