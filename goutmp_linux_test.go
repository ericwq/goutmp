// Copyright 2023~2024 wangqi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goutmp

import (
	"os/user"
	"testing"
)

func TestGetRecord(t *testing.T) {
	v := GetRecord()

	if v == nil {
		t.Skip("this system doesn't support utmpx access.")
	}

	c := 0
	for v != nil {
		t.Logf("type=%2d, pid=%6d, line=%7s, id=%2d, user=%10s, host=%29s, exit=%v, session=%d, time=%s\n",
			v.Type, v.GetPid(), v.GetLine(), v.GetId(), v.GetUser(),
			v.GetHost(), v.Exit, v.Session, v.GetTime())
		v = GetRecord()
		c++
	}
	if c == 0 {
		t.Errorf("#test GetUtmpx return nothing.")
	}
	if c == 0 {
		t.Errorf("#test GetUtmpx got %d records\n", c)
	}
}

func TestDeviceExists(t *testing.T) {
	tc := []struct {
		label  string
		line   string
		expect bool
	}{
		{"pts/0", "pts/0", true},
		{"tty", "tty", true},
		// {"tty0", "tty0", true}, // this one doesn't work for some linux container,
		{"tty/1", "tty/1", false},
	}

	for _, v := range tc {
		t.Run(v.label, func(t *testing.T) {
			got := DeviceExists(v.line)
			if got != v.expect {
				t.Errorf("#test %s expect %t, got %t\n", v.label, v.expect, got)
			}
		})
	}
}

func TestUtmpxAPI(t *testing.T) {
	// for normal user, they can't update utmp/wtmp
	v := AddRecord("/devpts/0", "unknow", "192.168.0.1", 12)
	if v {
		t.Error("#test AddRecord() failed")
	}

	v = RemoveRecord("/devpts/0", 12)
	if v {
		t.Error("#test RemoveRecord() failed")
	}

	u, _ := user.Current()
	v = AddLastLog("/dev/pts/0", u.Username, "127.0.0.1")
	if v {
		t.Error("#test PutLastlogEntry() failed")
	}
	v = AddLastLog("/dev/pts/0", "unknow", "127.0.0.1")
	if v {
		t.Error("#test PutLastlogEntry() failed")
	}
}
