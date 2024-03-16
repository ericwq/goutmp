// Copyright 2023~2024 wangqi. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package goutmp

func AddRecord(ptsName string, user string, host string, pid int) bool {
	// fmt.Fprintf(os.Stderr, "unimplement %s\n", "AddRecord()")
	return false
}

func RemoveRecord(ptsName string, pid int) bool {
	// fmt.Fprintf(os.Stderr, "unimplement %s\n", "RemoveRecord()")
	return false
}

func GetRecord() *Utmpx {
	// fmt.Fprintf(os.Stderr, "unimplement %s\n", "GetRecord()")
	return nil
}

func AddLastLog(line, userName, host string) bool {
	// fmt.Fprintf(os.Stderr, "unimplement %s\n", "PutLastlogEntry()")
	return false
}
