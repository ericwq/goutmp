//go:build linux

package main

import (
	// "time"

	utmps "github.com/ericwq/goutmp"
)

func main() {
	user := "ide"
	host := "192.168.1.10"
	utmp := utmps.Put_utmp(user, "/dev/pts/5", host)
	utmps.Put_lastlog_entry("xs", user, "/dev/pts/5", host)
	// time.Sleep(10 * time.Second)
	utmps.Unput_utmp(utmp)
}
