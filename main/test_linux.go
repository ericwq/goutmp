//go:build linux

package main

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/creack/pty"
	utmps "github.com/ericwq/goutmp"
)

// run with
// go run -tags utmps .
// go run -tags utmp  .
func main() {
	// user := "ide"
	// host := "192.168.1.10"
	// utmp := utmps.Put_utmp(user, "/dev/pts/5", host)
	// utmps.Put_lastlog_entry("xs", user, "/dev/pts/5", host)
	// time.Sleep(10 * time.Second)
	// utmps.Unput_utmp(utmp)

	host := "192.168.1.10"
	ptmx, pts, err := pty.Open()
	fmt.Printf("#test main pts=%s, ptmx=%s\n", pts.Name(), ptmx.Name())
	if err != nil {
		fmt.Printf("#test open pts error:%s\n", err)
	}
	u, _ := user.Current()
	v := utmps.GetRecord()
	if v == nil {
		fmt.Println("this system doesn't support utmpx access.")
	}

	if ok := utmps.AddRecord(pts.Name(), u.Username, host, os.Getpid()); !ok {
		fmt.Printf("#test UtmpxAddRecord retrun false\n")
	}

	// this one doesn't work on linux
	// utmps.PutLastlogEntry("line", "ide", pts.Name())

	time.Sleep(10 * time.Second)
	if ok := utmps.RemoveRecord(pts.Name(), os.Getpid()); !ok {
		fmt.Printf("#test UtmpxRemoveRecord return false\n")
	}
}
