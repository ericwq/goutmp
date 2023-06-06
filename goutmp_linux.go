package goutmp

/*
#cgo pkg-config: utmps skalibs

#include <pwd.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/file.h>
#include <time.h>
#include <unistd.h>

#include <utmps/utmps.h>
#include <lastlog.h>

// typedef char char_t;

void pututmp(struct utmpx* ut, char* uname, char* ptsname, char* host) {
	// printf("effective GID=%u\n", getegid());
	// system("echo ---- pre ----;who");
	memset(ut, 0, sizeof(struct utmpx));

	ut->ut_type = USER_PROCESS;  // This is a user login
	strncpy(ut->ut_user, uname, sizeof(ut->ut_user));
	time((time_t*)&(ut->ut_tv.tv_sec));  // Stamp with current time
	ut->ut_pid = getpid();

	// Set ut_line and ut_id based on the terminal associated with 'stdin'. This code assumes
	// terminals named "/dev/[pt]t[sy]*". The "/dev/" dirname is 5 characters; the "[pt]t[sy]"
	// filename prefix is 3 characters (making 8 characters in all).

	// devName = ttyname(STDIN_FILENO);
	// if (devName == NULL)
	// 	errExit("ttyname");
	// if (strlen(devName) <= 8) // Should never happen
	// 	fatal("Terminal name is too short: %s", devName);
	strncpy(ut->ut_line, ptsname + 5, sizeof(ut->ut_line));
	strncpy(ut->ut_id, ptsname + 8, sizeof(ut->ut_id));
	strcpy(ut->ut_host, host);

	setutxent();               // Rewind to start of utmp file
	pututxline(ut);            // Overwrite previous utmp record
	updwtmpx(_PATH_WTMP, ut);  // Append login record to wtmp
	endutxent();
	// system("echo ---- post ----;who");
}

void unpututmp(struct utmpx* ut) {
	ut->ut_type = DEAD_PROCESS;              // Required for logout record
	time((time_t*)&(ut->ut_tv.tv_sec));      // Stamp with logout time
	memset(&(ut->ut_user), 0, UT_NAMESIZE);  // Logout record has null username
	setutxent();
	pututxline(ut);
	updwtmpx(_PATH_WTMP, ut);  // Append logout record to wtmp
	endutxent();

	// system("echo ---- cleanup ----;who; last");
}

struct utmpx* res = NULL;

struct utmpx* getutmp() {
	if (res != NULL)  // If 'res' was set via a previous call
		memset(res, 0, sizeof(struct utmpx));
	res = getutxent();
	if (res == NULL) {
		return NULL;
	}

	// unsigned char* charPtr = (unsigned char*)res;
	// int i;
	// int start = 32+4+4+2;
	// int end = start+4;  // sizeof(struct utmpx); )
	// for (i = start; i < end; i++)
	// 	printf("%02x ", charPtr[i]);
	// printf("\n");

	printf(
		"[ C] type=%d; pid=%d; line=%s, id=%.4s; user=%s; host=%s; exit={%u %u}, session=%d "
		"time={%ld %ld}\n",
		res->ut_type, res->ut_pid, res->ut_line, res->ut_id, res->ut_user, res->ut_host,
		res->ut_exit.e_termination, res->ut_exit.e_exit, res->ut_session, res->ut_tv.tv_sec,
		res->ut_tv.tv_usec);
	return res;
}

int putlastlogentry(int64_t t, int uid, char* line, char* host) {
	int retval = 0;
	FILE* f;
	struct lastlog l;

	strncpy(l.ll_line, line, UT_LINESIZE);
	l.ll_line[UT_LINESIZE - 1] = '\0';
	strncpy(l.ll_host, host, UT_HOSTSIZE);
	l.ll_host[UT_HOSTSIZE - 1] = '\0';

	l.ll_time = (time_t)t;
	// printf("l: ll_line '%s', ll_host '%s', ll_time %d\n", l.ll_line, l.ll_host, l.ll_time);

	// Write lastlog entry at fixed offset (uid * sizeof(struct lastlog)
	if (NULL != (f = fopen("/var/log/lastlog", "rw+"))) {
		if (!fseek(f, (uid * sizeof(struct lastlog)), SEEK_SET)) {
			int fd = fileno(f);
			if (write(fd, &l, sizeof(l)) == sizeof(l)) {
				retval = 1;
				// int32_t stat = system("echo ---- lastlog ----; lastlog");
			}
		}
		fclose(f);
	}
	return retval;
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os/user"
	"strings"
	"time"
	"unsafe"
)

// UtmpEntry wraps the C struct utmp
type UtmpEntry struct {
	entry C.struct_utmpx
}

// func (u *Utmpx) GetLine() string {
// 	return unsafe.Slice(u.Line,32)
// 	return fmt.Sprintf("%s", u.Line)
// }

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

var hostEndian binary.ByteOrder

func init() {
	// https://commandcenter.blogspot.com/2012/04/byte-order-fallacy.html
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		hostEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		hostEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

// read the next utmpx record from utmp database
func GetUtmpx() *Utmpx {
	/*
		https://github.com/llgoer/cgo-struct-array/blob/master/src/main.go
		https://medium.com/@liamkelly17/working-with-packed-c-structs-in-cgo-224a0a3b708b
		https://github.com/brgl/busybox/blob/master/coreutils/who.c
	*/
	g := &Utmpx{}

	p := C.getutmp()
	// p := C.getutxent()
	if p == nil {
		return nil
	}
	// convert C struct into Go struct for utmpx
	cdata := C.GoBytes(unsafe.Pointer(p), C.sizeof_struct_utmpx)
	buf := bytes.NewBuffer(cdata)

	binary.Read(buf, hostEndian, g)

	// convert C struct into Go struct for exit_status
	data2 := C.GoBytes(unsafe.Pointer(&p.ut_exit), C.sizeof_struct_exit_status)
	buf2 := bytes.NewBuffer(data2)
	s2 := &ExitStatus{}
	binary.Read(buf2, hostEndian, &s2.Termination)
	binary.Read(buf2, hostEndian, &s2.Exit)
	g.Exit = *s2

	// convert C struct into Go struct for timeval
	data3 := C.GoBytes(unsafe.Pointer(&p.ut_tv), C.sizeof_struct_timeval)
	buf3 := bytes.NewBuffer(data3)
	s3 := &TimeVal{}
	binary.Read(buf3, hostEndian, &s3.Sec)
	binary.Read(buf3, hostEndian, &s3.Usec)
	g.Tv = *s3

	data2 = C.GoBytes(unsafe.Pointer(&p.ut_pid), C.sizeof_pid_t)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, binary.LittleEndian, &(g.Pid))

	data2 = C.GoBytes(unsafe.Pointer(&p.ut_id), UTMPS_UT_IDSIZE)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, binary.LittleEndian, &(g.Id))

	data2 = C.GoBytes(unsafe.Pointer(&p.ut_user), UTMPS_UT_NAMESIZE)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, binary.LittleEndian, &(g.User))

	data2 = C.GoBytes(unsafe.Pointer(&p.ut_session), C.sizeof_pid_t)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, binary.LittleEndian, &(g.Session))
	return g
}

func (u *Utmpx) GetType() int {
	return int(u.Type)
}

func (u *Utmpx) GetPid() int {
	return int(u.Pid)
}

func (u *Utmpx) GetUser() string {
	return B2S(u.User[:])
}

func (u *Utmpx) GetHost() string {
	return B2S(u.Host[:])
}

func (u *Utmpx) GetLine() string {
	return B2S(u.Line[:])
}

func (u *Utmpx) GetId() string {
	return B2S(u.Id[:4])
}

func (u *Utmpx) GetTime() time.Time {
	return time.Unix(u.Tv.Sec, u.Tv.Usec)
}

// convert int8 arrary to string
func B2S(bs []int8) string {
	//	https://stackoverflow.com/questions/28848187/how-to-convert-int8-to-string

	ba := make([]byte, 0, len(bs))
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
