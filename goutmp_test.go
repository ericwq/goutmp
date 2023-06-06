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
		fmt.Printf("[Go] USER:%s, HOST:%s, ID=%s, LINE=%s, TIME=%v \n", v.GetUser(), v.GetHost(), v.GetId(), v.GetLine(), v.Tv)
		v = GetUtmpx()
	}
	if v != nil {
		t.Errorf("#test GetUtmpx should return nil now.")
	}
}
