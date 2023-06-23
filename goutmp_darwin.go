package goutmp

import "os"

func UtmpxAddRecord(pts *os.File, host string) bool {
	// not implement
	return false
}

func UtmpxRemoveRecord(pts *os.File) bool {
	// not implement
	return false
}

func GetUtmpx() *Utmpx {
	// not implement
	return nil
}

func PutLastlogEntry(line, userName, host string) bool {
	// not implement
	return false
}
