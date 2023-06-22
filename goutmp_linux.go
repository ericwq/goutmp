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

#define MIN(a_, b_) (((a_) < (b_)) ? (a_) : (b_))

static int write_uwtmp_record(const char* user,
							  const char* termName,
							  const char* host,
							  pid_t pid,
							  int add) {
	struct utmpx ut;
	struct timeval tv;
	size_t len, offset;

	memset(&ut, 0, sizeof(ut));

	memset(&tv, 0, sizeof(tv));
	(void)gettimeofday(&tv, 0);

	if (user) {  // Logout record has null username
		len = strlen(user);
		memcpy(ut.ut_user, user, MIN(sizeof(ut.ut_user), len));
	}

	if (host) {  // Logout record has null host
		len = strlen(host);
		memcpy(ut.ut_host, host, MIN(sizeof(ut.ut_host), len));
	}

	// len = strlen(term);
	// memcpy(ut.ut_line, term, MIN(sizeof(ut.ut_line), len));
	//
	// offset = len <= sizeof(ut.ut_id) ? 0 : len - sizeof(ut.ut_id);
	// memcpy(ut.ut_id, term + offset, len - offset);

	// Set ut_line and ut_id based on the terminal associated with 'stdin'. This
	// code assumes terminals named "/dev/[pt]t[sy]*". The "/dev/" dirname is 5
	// characters; the "[pt]t[sy]" filename prefix is 3 characters (making 8
	// characters in all).

	len = strlen(termName + 5);
	memcpy(ut.ut_line, termName + 5, MIN(sizeof(ut.ut_line), len));

	len = strlen(termName + 8);
	memcpy(ut.ut_id, termName + 8, MIN(sizeof(ut.ut_id), len));

	if (add)
		ut.ut_type = USER_PROCESS;
	else
		ut.ut_type = DEAD_PROCESS;

	ut.ut_pid = pid;

	ut.ut_tv.tv_sec = (__typeof__(ut.ut_tv.tv_sec))tv.tv_sec;
	ut.ut_tv.tv_usec = (__typeof__(ut.ut_tv.tv_usec))tv.tv_usec;

	setutxent();
	if (!pututxline(&ut))
		return EXIT_FAILURE;
	// fatal_error("pututline: %s", strerror(errno));
	endutxent();

	(void)updwtmp(_PATH_WTMP, &ut);

	// debug_msg("utmp/wtmp record %s for terminal '%s'",
	// 	  add ? "added" : "removed", term);
	return EXIT_SUCCESS;
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
	//
	// printf(
	// 	"[ C] type=%d; pid=%d; line=%s, id=%.4s; user=%s; host=%s; exit={%u %u}; session=%d; "
	// 	"time={%ld %ld}\n",
	// 	res->ut_type, res->ut_pid, res->ut_line, res->ut_id, res->ut_user, res->ut_host,
	// 	res->ut_exit.e_termination, res->ut_exit.e_exit, res->ut_session, res->ut_tv.tv_sec,
	// 	res->ut_tv.tv_usec);
	return res;
}

// return 1 means success, otherwise return 0.
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
	"os"
	"os/user"
	"strings"
	"time"
	"unsafe"
)

// UtmpEntry wraps the C struct utmp
// type UtmpEntry struct {
// 	entry C.struct_utmpx
// }

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

/*
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
*/

// adds a login record to the database for the TTY belonging to
// the pseudo-terminal slave file pts, using the username corresponding with the
// real user ID of the calling process and the optional hostname host.
func UtmpxAddRecord(pts *os.File, host string) bool {
	user, err := user.Current()
	if err != nil {
		return false
	}

	termName := C.CString(pts.Name())
	userName := C.CString(user.Username)
	hostName := C.CString(host)
	pid := os.Getpid()
	defer func() {
		C.free(unsafe.Pointer(termName))
		C.free(unsafe.Pointer(userName))
		C.free(unsafe.Pointer(hostName))
	}()

	// C.pututmp(&entry, userName, ptsName, hostName)
	return C.write_uwtmp_record(userName, termName, hostName, C.pid_t(pid), 1) == 0
}

// marks the login session as being closed for the TTY belonging to the
// pseudo-terminal slave file pts, using the PID of the calling process
func UtmpxRemoveRecord(pts *os.File) bool {
	// git clone https://git.launchpad.net/ubuntu/+source/libutempter

	termName := C.CString(pts.Name())
	pid := os.Getpid()
	defer func() {
		C.free(unsafe.Pointer(termName))
	}()

	// ut_type, ut_id and ut_line is required, ut_user must be zero
	return C.write_uwtmp_record(nil, termName, nil, C.pid_t(pid), 1) == 0
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

	// convert exit field
	data2 := C.GoBytes(unsafe.Pointer(&p.ut_exit), C.sizeof_struct_exit_status)
	buf2 := bytes.NewBuffer(data2)
	s2 := &ExitStatus{}
	binary.Read(buf2, hostEndian, &s2.Termination)
	binary.Read(buf2, hostEndian, &s2.Exit)
	g.Exit = *s2

	// convert tv field
	data3 := C.GoBytes(unsafe.Pointer(&p.ut_tv), C.sizeof_struct_timeval)
	buf3 := bytes.NewBuffer(data3)
	s3 := &TimeVal{}
	binary.Read(buf3, hostEndian, &s3.Sec)
	binary.Read(buf3, hostEndian, &s3.Usec)
	g.Tv = *s3

	// convert pid field
	data2 = C.GoBytes(unsafe.Pointer(&p.ut_pid), C.sizeof_pid_t)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, hostEndian, &(g.Pid))

	// convert id field
	data2 = C.GoBytes(unsafe.Pointer(&p.ut_id), UTMPS_UT_IDSIZE)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, hostEndian, &(g.Id))

	// convert user field
	data2 = C.GoBytes(unsafe.Pointer(&p.ut_user), UTMPS_UT_NAMESIZE)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, hostEndian, &(g.User))

	// convert session field
	data2 = C.GoBytes(unsafe.Pointer(&p.ut_session), C.sizeof_pid_t)
	buf2 = bytes.NewBuffer(data2)
	binary.Read(buf2, hostEndian, &(g.Session))
	return g
}

func (u *Utmpx) GetType() int {
	return int(u.Type)
}

func (u *Utmpx) GetPid() int {
	return int(u.Pid)
}

func (u *Utmpx) GetUser() string {
	return b2s(u.User[:UTMPS_UT_NAMESIZE])
}

func (u *Utmpx) GetHost() string {
	return b2s(u.Host[:UTMPS_UT_HOSTSIZE])
}

func (u *Utmpx) GetLine() string {
	return b2s(u.Line[:UTMPS_UT_LINESIZE])
}

func (u *Utmpx) GetId() string {
	return b2s(u.Id[:UTMPS_UT_IDSIZE])
}

func (u *Utmpx) GetTime() time.Time {
	return time.Unix(u.Tv.Sec, u.Tv.Usec)
}

// convert int8 arrary to string
func b2s(bs []int8) string {
	//	https://stackoverflow.com/questions/28848187/how-to-convert-int8-to-string

	ba := make([]byte, 0, len(bs))
	for _, b := range bs {
		if b == 0 { // skip zero
			continue
		}
		ba = append(ba, byte(b))
	}
	return string(ba)
}

// Put the login line, username and originating host/IP to lastlog
func PutLastlogEntry(line, userName, host string) bool {
	u, e := user.Lookup(userName)
	if e != nil {
		return false
	}
	var uid uint32
	fmt.Sscanf(u.Uid, "%d", &uid)

	t := time.Now().Unix()
	lineC := C.CString(line)
	hostC := C.CString(host)
	defer func() {
		C.free(unsafe.Pointer(lineC))
		C.free(unsafe.Pointer(hostC))
	}()

	return C.putlastlogentry(C.int64_t(t), C.int(uid), lineC, hostC) == 1
	// stat := C.putlastlogentry(C.int64_t(t), C.int(uid), C.CString(app), C.CString(host))
	// fmt.Println("stat was:",stat)
}

// return true if we can read from the utmp data file
func HasUtmpSupport() bool {
	r := GetUtmpx()
	if r != nil {
		return true
	}
	return false
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
