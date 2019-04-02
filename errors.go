package pid

import (
	"fmt"
)

var (
	// Segregated to simplify localization
	errStalePid    = fmt.Errorf("stale pid check open failed")
	errFileLocked  = fmt.Errorf("pid file locked by running process")
	errOpenFailed  = fmt.Errorf("failed to open pid file")
	errWriteFailed = fmt.Errorf("failed to write pid to file")
	errShortWrite  = fmt.Errorf("short write failure writing pid")
	errCleanFailed = fmt.Errorf("failed to clean up pid file")
)
