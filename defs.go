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
        _T_EMPTY         = C.EMPTY
        _T_USER_PROCESS  = C.USER_PROCESS
        _T_INIT_PROCESS  = C.INIT_PROCESS
        _T_LOGIN_PROCESS = C.LOGIN_PROCESS
        _T_DEAD_PROCESS  = C.DEAD_PROCESS
)
