package goutmp

import (
	"fmt"
	"os"
)

// check whether device exist, the line parameter should be form of "pts/0" or "ttys001"
func DeviceExists(line string) bool {
	deviceName := fmt.Sprintf("/dev/%s", line)
	_, err := os.Lstat(deviceName)
	if err != nil {
		return false
	}

	return true
}
