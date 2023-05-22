#include <pwd.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/file.h>
#include <time.h>
#include <unistd.h>

#include <lastlog.h>
#include <utmpx.h>

// #cgo CFLAGS: -I/usr/include/utmps -v
// #cgo LDFLAGS: -L/lib -lutmps -lskarnet

// #include <lastlog.h>
// gcc -I/usr/include/utmps -v -L/lib -lutmps -lskarnet -c -o xutmp.o xutmp.c
void pututmp(struct utmpx* entry, char* uname, char* ptsname, char* host);
void unpututmp(struct utmpx* entry);
int putlastlogentry(int64_t t, int uid, char* line, char* host);
