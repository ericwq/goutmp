#include <pwd.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/file.h>
#include <time.h>
#include <unistd.h>

#include <utmps/utmps.h>

void pututmp(struct utmpx* entry, char* uname, char* ptsname, char* host);
void unpututmp(struct utmpx* entry);
int putlastlogentry(int64_t t, int uid, char* line, char* host);
