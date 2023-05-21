// +build freebsd

package main

import (
	"time"

	"github.com/ericwq/goutmp"
)

func main() {
	user := "bin"
	host := "test.example.com"
	utmp := goutmp.Put_utmp(user, "/dev/pts0", host)
	goutmp.Put_lastlog_entry("xs", user, "/dev/pts0", host)
	time.Sleep(10 * time.Second)
	goutmp.Unput_utmp(utmp)
}
