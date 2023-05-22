package goutmp

//#cgo CFLAGS: -I/usr/include/utmps -v
//#cgo LDFLAGS: -L/lib -lutmps -lskarnet
//#include <stdio.h>
//#include <stdlib.h>
//#include <sys/file.h>
//#include <string.h>
//#include <unistd.h>
//#include <stdint.h>
//#include <time.h>
//#include <pwd.h>
//
//#include <utmps/utmps.h>
//#include <lastlog.h>
//
//typedef char char_t;
//
//
//void pututmp(struct utmpx* entry, char* uname, char* ptsname, char* host) {
//  //int32_t stat = system("echo ---- pre ----;who");
//
//  entry->ut_type = USER_PROCESS;
//  entry->ut_pid = getpid();
//  strcpy(entry->ut_line, ptsname + strlen("/dev/"));
//  /* only correct for ptys named /dev/tty[pqr][0-9a-z] */
//
//  strcpy(entry->ut_id, ptsname + strlen("/dev/pts/"));
//  entry->ut_time = time(NULL);
//  strcpy(entry->ut_user, uname);
//  strcpy(entry->ut_host, host);
//  entry->ut_addr = 0;
//  setutxent();
//  pututxline(entry);
//
//  system("echo ---- post ----;who");
//  printf("line:%s, id:%s, user:%s, host:%s\n", entry->ut_line, entry->ut_id, entry->ut_user, entry->ut_host);
//}
//
//void unpututmp(struct utmpx* entry) {
//  entry->ut_type = DEAD_PROCESS;
//  memset(entry->ut_line, 0, UT_LINESIZE);
//  entry->ut_time = 0;
//  memset(entry->ut_user, 0, UT_NAMESIZE);
//  setutxent();
//  pututxline(entry);
//
//  system("echo ---- cleanup ----;who; last");
//
//  endutxent();
//}
//
//int putlastlogentry(int64_t t, int uid, char* line, char* host) {
//  int retval = 0;
//  FILE *f;
//  struct lastlog l;
//
//  strncpy(l.ll_line, line, UT_LINESIZE);
//  l.ll_line[UT_LINESIZE-1] = '\0';
//  strncpy(l.ll_host, host, UT_HOSTSIZE);
//  l.ll_host[UT_HOSTSIZE-1] = '\0';
//
//  l.ll_time = (time_t)t;
//  //printf("l: ll_line '%s', ll_host '%s', ll_time %d\n", l.ll_line, l.ll_host, l.ll_time);
//
//  /* Write lastlog entry at fixed offset (uid * sizeof(struct lastlog) */
//  if( NULL != (f = fopen("/var/log/lastlog", "rw+")) ) {
//    if( !fseek(f, (uid * sizeof(struct lastlog)), SEEK_SET) ) {
//      int fd = fileno(f);
//      if( write(fd, &l, sizeof(l)) == sizeof(l) ) {
//        retval = 1;
//        //int32_t stat = system("echo ---- lastlog ----; lastlog");
//      }
//    }
//    fclose(f);
//  }
//  return retval;
//}
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
