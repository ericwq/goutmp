package goutmp

/*
#cgo CFLAGS: -I./xutmp
#cgo LDFLAGS: -L${SRCDIR}/xutmp -lxutmp

#include "xutmp.h"
*/
import "C"

import (
	"fmt"
	"net"
	"os/user"
	"strings"
	"time"
)

// UtmpEntry wraps the C struct utmp
type UtmpEntry struct {
	entry C.struct_utmpx
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

// Put a username and the originating host/IP to utmp
func Put_utmp(user, ptsName, host string) UtmpEntry {
	var entry UtmpEntry

	// log.Println("Put_utmp:host ", host, " user ", user)
	C.pututmp(&entry.entry, C.CString(user), C.CString(ptsName), C.CString(host))
	return entry
}

// Remove a username/host entry from utmp
func Unput_utmp(entry UtmpEntry) {
	C.unpututmp(&entry.entry)
}

// Put the login app, username and originating host/IP to lastlog
func Put_lastlog_entry(app, usr, ptsname, host string) {
	u, e := user.Lookup(usr)
	if e != nil {
		return
	}
	var uid uint32
	fmt.Sscanf(u.Uid, "%d", &uid)

	t := time.Now().Unix()
	_ = C.putlastlogentry(C.int64_t(t), C.int(uid), C.CString(app), C.CString(host))
	// stat := C.putlastlogentry(C.int64_t(t), C.int(uid), C.CString(app), C.CString(host))
	// fmt.Println("stat was:",stat)
}
