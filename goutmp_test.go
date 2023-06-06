package goutmp

import (
	"fmt"
	"testing"
)

func TestGetUtmpx(t *testing.T) {
	v := GetUtmpx()

	if v == nil {
		t.Errorf("#test failed to get utmp record. %v", v)
	}
	for v != nil {
		fmt.Printf("[Go] type=%d, pid=0x%8x, line=%s, id=%s, user=%s, host=%s, time=%v \n",
			v.Type, v.GetPid(), v.GetLine(), v.GetId(), v.GetUser(), v.GetHost(), v.Tv)
		v = GetUtmpx()
	}
	if v != nil {
		t.Errorf("#test GetUtmpx should return nil now.")
	}
}
