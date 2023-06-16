//go:build linux

package main

import (
	"fmt"
	"time"

	"github.com/creack/pty"
	utmps "github.com/ericwq/goutmp"
)

func main() {
	// user := "ide"
	// host := "192.168.1.10"
	// utmp := utmps.Put_utmp(user, "/dev/pts/5", host)
	// utmps.Put_lastlog_entry("xs", user, "/dev/pts/5", host)
	// time.Sleep(10 * time.Second)
	// utmps.Unput_utmp(utmp)

	// user := "ide"
	host := "192.168.1.10"
	_, pts, err := pty.Open()
	if err != nil {
		fmt.Printf("#test open pts error:%s\n", err)
	}
	if ok := utmps.UtmpxAddRecord(pts, host); !ok {
		fmt.Printf("#test UtmpxAddRecord retrun false\n")
	}
	// utmps.Put_lastlog_entry("xs", user, "/dev/pts/5", host)
	time.Sleep(10 * time.Second)
	if ok := utmps.UtmpxRemoveRecord(pts); !ok {
		fmt.Printf("#test UtmpxRemoveRecord return false\n")
	}
}
