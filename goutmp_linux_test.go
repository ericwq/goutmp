package goutmp

import (
	// "fmt"
	"testing"
)

func TestGetUtmpx(t *testing.T) {
	v := GetUtmpx()

	if v == nil {
		t.Errorf("#test failed to get utmp record. %v", v)
	}

	c := 0
	// fmt.Printf("20018=%x\n", 20018)
	for v != nil {
		// fmt.Printf("[Go] type=%d, pid=%d, line=%s, id=%s, user=%s, host=%s, exit=%v, session=%d, time=%v\nt=%s\n",
		// 	v.Type, v.GetPid(), v.GetLine(), v.GetId(), v.GetUser(), v.GetHost(), v.Exit, v.Session, v.Tv, v.GetTime())
		v = GetUtmpx()
		c++
	}
	if v != nil {
		t.Errorf("#test GetUtmpx should return nil now.")
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
		{"pts/1", "pts/1", true},
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
