package goutmp

import (
	// "fmt"
	"testing"
)

func TestGetUtmpxDarwin(t *testing.T) {
	v := GetUtmpx()

	if v != nil {
		t.Errorf("#test expect nil, get utmp record. %v", v)
	}
}

func TestDeviceExistsDarwin(t *testing.T) {
	tc := []struct {
		label  string
		line   string
		expect bool
	}{
		{"ttys001", "ttys001", true},
		{"ttys1", "ttys1", true},
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

