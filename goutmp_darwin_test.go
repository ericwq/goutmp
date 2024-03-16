// Copyright 2023~2024 wangqi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goutmp

import (
	// "fmt"

	"testing"
)

func TestDeviceExists(t *testing.T) {
	tc := []struct {
		label  string
		line   string
		expect bool
	}{
		{"ttys001", "ttys001", true},
		{"ttys1", "ttys1", true},
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
	tc := []struct {
		label  string
		expect bool
	}{
		{"AddRecord", false},
		{"RemoveRecord", false},
		{"GetRecord", false},
		{"PutLastlogEntry", false},
	}

	for _, v := range tc {
		t.Run(v.label, func(t *testing.T) {

			switch v.label {
			case "AddRecord":
				ret := AddRecord("", "", "", 0)
				if ret {
					t.Errorf("AddRecord is not implemented")
				}
			case "RemoveRecord":
				ret := RemoveRecord("", 0)
				if ret {
					t.Errorf("RemoveRecord is not implemented")
				}
			case "GetRecord":
				ret := GetRecord()
				if ret != nil {
					t.Errorf("GetRecord is not implemented")
				}
			case "PutLastlogEntry":
				ret := AddLastLog("", "", "")
				if ret {
					t.Errorf("PutLastlogEntry is not implemented")
				}
			}

		})
	}
}
