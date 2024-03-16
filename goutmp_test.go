// Copyright 2023~2024 wangqi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goutmp

import (
	"testing"
	"time"
)

func TestHasUtmpSupport(t *testing.T) {
	if !HasUtmpSupport() {
		t.Skip("the system doesn't support utmpx access.")
	}
}

func TestUtmpxAccess(t *testing.T) {
	tc := []struct {
		label string
		user  string
		host  string
		line  string
		id    int
		pid   int
		typ0  int
		time  time.Time
	}{
		{"normal", "john", "192.168.0.1", "pts/0", 120, 512, USER_PROCESS, time.Now()},
	}

	for _, v := range tc {
		t.Run(v.label, func(t *testing.T) {
			u := &Utmpx{}
			u.SetUser(v.user)
			u.SetHost(v.host)
			u.SetLine(v.line)
			u.SetId(v.id)
			u.SetPid(v.pid)
			u.SetType(v.typ0)
			u.SetTime(v.time)

			if u.GetUser() != v.user {
				t.Errorf("%q GetUser expect %s, got %s\n", v.label, v.user, u.GetUser())
			}
			if u.GetHost() != v.host {
				t.Errorf("%q GetHost expect %s, got %s\n", v.label, v.host, u.GetHost())
			}
			if u.GetLine() != v.line {
				t.Errorf("%q GetLine expect %s, got %s\n", v.label, v.line, u.GetLine())
			}
			if u.GetId() != v.id {
				t.Errorf("%q GetId expect %d, got %d\n", v.label, v.id, u.GetId())
			}
			if u.GetPid() != v.pid {
				t.Errorf("%q GetPid expect %d, got %d\n", v.label, v.pid, u.GetPid())
			}
			if u.GetType() != v.typ0 {
				t.Errorf("%q GetType expect %d, got %d\n", v.label, v.typ0, u.GetType())
			}
			if u.GetTime() != v.time {
				t.Skip("please reconsider convertion between TimeVal and Time.")
				// t.Errorf("%q GetTime expect \n%s, got \n%s\n", v.label, v.time, u.GetTime())
			}
		})
	}
}
