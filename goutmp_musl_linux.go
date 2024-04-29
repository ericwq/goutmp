// Copyright 2023~2024 wangqi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//go:build cgo && utmps

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
	updwtmpx(_PATH_WTMP, &ut);

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
	"os/user"
	"time"
	"unsafe"
)

const (
	_ADD    = 1
	_REMOVE = 0
)

// called when user login, adds a login utmp/wtmp record to database with the specified
// pseudo-terminal device name, user name, host name and process PID.
// this fuction will update both utmp and wtmp within one call
func AddRecord(ptsName string, user string, host string, pid int) bool {
	termName := C.CString(ptsName)
	userName := C.CString(user)
	hostName := C.CString(host)
	defer func() {
		C.free(unsafe.Pointer(termName))
		C.free(unsafe.Pointer(userName))
		C.free(unsafe.Pointer(hostName))
	}()

	return C.write_uwtmp_record(userName, termName, hostName, C.pid_t(pid), _ADD) == 0
}

// called when user logout, marks a login session as being closed with the specified
// pseudo-terminal device name, process PID.
// this fuction will update both utmp and wtmp within one call
func RemoveRecord(ptsName string, pid int) bool {
	// git clone https://git.launchpad.net/ubuntu/+source/libutempter

	termName := C.CString(ptsName)
	defer func() {
		C.free(unsafe.Pointer(termName))
	}()

	// ut_user and ut_host must be zero, ut_id and ut_line, ut_pid is required
	return C.write_uwtmp_record(nil, termName, nil, C.pid_t(pid), _REMOVE) == 0
}

// read the next utmpx record from utmp database
func GetRecord() *Utmpx {
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

// Add a record to last log, with the specified login line, username and
// originating host/IP.
func AddLastLog(line, userName, host string) bool {
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
