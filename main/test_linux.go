//go:build linux

package main

import (
	// "time"

	utmps "github.com/ericwq/goutmp"
)

func main() {
	user := "ide"
	host := "example.com"
	utmp := utmps.Put_utmp(user, "/dev/pts/0", host)
	utmps.Put_lastlog_entry("xs", user, "/dev/pts/0", host)
	// time.Sleep(10 * time.Second)
	utmps.Unput_utmp(utmp)
}
