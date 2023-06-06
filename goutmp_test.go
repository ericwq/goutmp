package goutmp

import (
	"testing"
)

func TestGetUtmpx(t *testing.T) {
	v := GetUtmpx()

	if v == nil {
		t.Errorf("#test failed to get utmp record. %v", v)
	}

	c:=0
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
	if c==0 {
		t.Errorf("#test GetUtmpx got %d records\n", c)
	}
}
