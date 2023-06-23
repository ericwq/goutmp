package goutmp

import (
	"fmt"
	"net"
	"os"
	"strings"
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

// return true if we can read from the utmp data file
func HasUtmpSupport() bool {
	r := GetUtmpx()
	if r != nil {
		return true
	}
	return false
}

// return remote client hostname or IP if host lookup fails
// addr is expected to be of the format given by net.Addr.String()
// eg., "127.0.0.1:80" or "[::1]:80"
func GetHost(addr string) (h string) {
	if !strings.Contains(addr, "[") {
		h = strings.Split(addr, ":")[0]
	} else {
		h = strings.Split(strings.Split(addr, "[")[1], "]")[0]
	}
	hList, e := net.LookupAddr(h)
	// fmt.Printf("lookupAddr:%v\n", hList)
	if e == nil {
		h = hList[0]
	}
	return
}
