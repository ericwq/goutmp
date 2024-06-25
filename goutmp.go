// Copyright 2023~2024 wangqi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goutmp

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"
)

// set GOFLAGS="-tags=utmps" before nvim for musl based linux
// set GOFLAGS="-tags=utmp" before nvim for glibc based linux

// check whether device exist, the line parameter should be form of "pts/0" or "ttys001"
func DeviceExists(line string) bool {
	deviceName := fmt.Sprintf("/dev/%s", line)
	_, err := os.Lstat(deviceName)
	return err == nil
}

// return true if we can read from the utmp data file
func HasUtmpSupport() bool {
	r := GetRecord()
	return r != nil
}

// // return remote client hostname or IP if host lookup fails
// // addr is expected to be of the format given by net.Addr.String()
// // eg., "127.0.0.1:80" or "[::1]:80"
// func GetHost(addr string) (h string) {
// 	if !strings.Contains(addr, "[") {
// 		h = strings.Split(addr, ":")[0]
// 	} else {
// 		h = strings.Split(strings.Split(addr, "[")[1], "]")[0]
// 	}
// 	hList, e := net.LookupAddr(h)
// 	// fmt.Printf("lookupAddr:%v\n", hList)
// 	if e == nil {
// 		h = hList[0]
// 	}
// 	return
// }

func (u *Utmpx) GetHost() string { return b2s(u.Host[:UTMPS_UT_HOSTSIZE]) }

func (u *Utmpx) GetId() int {
	i, _ := strconv.Atoi(b2s(u.Id[:UTMPS_UT_IDSIZE]))
	return i
}
func (u *Utmpx) GetLine() string    { return b2s(u.Line[:UTMPS_UT_LINESIZE]) }
func (u *Utmpx) GetPid() int        { return int(u.Pid) }
func (u *Utmpx) GetTime() time.Time { return time.Unix(u.Tv.Sec, u.Tv.Usec) }
func (u *Utmpx) GetType() int       { return int(u.Type) }
func (u *Utmpx) GetUser() string    { return b2s(u.User[:UTMPS_UT_NAMESIZE]) }

func (u *Utmpx) SetHost(s string) {
	data := []byte(s)
	for i := range u.Host {
		if i < len(data) {
			u.Host[i] = int8(data[i])
		} else {
			break
		}
	}
}

func (u *Utmpx) SetId(id int) {
	data := []byte(fmt.Sprintf("%d", id))

	for i := range u.Id {
		if i < len(data) && i < UTMPS_UT_IDSIZE {
			u.Id[i] = int8(data[i])
		} else {
			break
		}
	}
}

func (u *Utmpx) SetLine(s string) {
	data := []byte(s)
	for i := range u.Line {
		if i < len(data) {
			u.Line[i] = int8(data[i])
		} else {
			break
		}
	}
}
func (u *Utmpx) SetPid(pid int) { u.Pid = int32(pid) }

func (u *Utmpx) SetTime(t time.Time) {
	u.Tv.Sec = t.Unix()
	u.Tv.Usec = (t.UnixNano() / 1e3 % 1e3)
}
func (u *Utmpx) SetType(t int) { u.Type = int16(t) }

func (u *Utmpx) SetUser(s string) {
	data := []byte(s)
	for i := range u.User {
		if i < len(data) {
			u.User[i] = int8(data[i])
		} else {
			break
		}
	}
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

// return true if we can read from the utmp data file
// func HasUtmpSupport() bool {
// 	r := GetUtmpx()
// 	if r != nil {
// 		return true
// 	}
// 	return false
// }

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
