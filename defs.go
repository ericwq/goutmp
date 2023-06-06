//go:build ignore

package goutmp

// #include <utmps/utmpx.h>
import "C"

type (
	TimeVal    C.struct_timeval
	ExitStatus C.struct_exit_status
	Utmpx      C.struct_utmpx
)

const (
	EMPTY         = C.EMPTY
	BOOT_TIME     = C.BOOT_TIME
	OLD_TIME      = C.OLD_TIME
	NEW_TIME      = C.NEW_TIME
	USER_PROCESS  = C.USER_PROCESS
	INIT_PROCESS  = C.INIT_PROCESS
	LOGIN_PROCESS = C.LOGIN_PROCESS
	DEAD_PROCESS  = C.DEAD_PROCESS

	UTMPS_UT_LINESIZE = C.UTMPS_UT_LINESIZE
	UTMPS_UT_NAMESIZE = C.UTMPS_UT_NAMESIZE
	UTMPS_UT_HOSTSIZE = C.UTMPS_UT_HOSTSIZE
	UTMPS_UT_IDSIZE   = C.UTMPS_UT_IDSIZE
)
