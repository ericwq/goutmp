#include "xutmp.h"
#include <lastlog.h>

typedef char char_t;

void pututmp(struct utmpx* entry, char* uname, char* ptsname, char* host) {
	// int32_t stat = system("echo ---- pre ----;who");

	entry->ut_type = USER_PROCESS;
	entry->ut_pid = getpid();
	strcpy(entry->ut_line, ptsname + strlen("/dev/"));
	/* only correct for ptys named /dev/tty[pqr][0-9a-z] */

	strcpy(entry->ut_id, ptsname + strlen("/dev/pts/"));
	entry->ut_time = time(NULL);
	strcpy(entry->ut_user, uname);
	strcpy(entry->ut_host, host);
	entry->ut_addr = 0;
	setutxent();
	pututxline(entry);
	endutxent();

	system("echo ---- post ----;who");
	// printf("line:%s, id:%s, user:%s, host:%s\n", entry->ut_line, entry->ut_id, entry->ut_user,
	// 	   entry->ut_host);
}

void unpututmp(struct utmpx* entry) {
	entry->ut_type = DEAD_PROCESS;
	memset(entry->ut_line, 0, UT_LINESIZE);
	entry->ut_time = 0;
	memset(entry->ut_user, 0, UT_NAMESIZE);
	setutxent();
	pututxline(entry);
	endutxent();

	system("echo ---- cleanup ----;who; last");
}

int putlastlogentry(int64_t t, int uid, char* line, char* host) {
	int retval = 0;
	FILE* f;
	struct lastlog l;

	strncpy(l.ll_line, line, UT_LINESIZE);
	l.ll_line[UT_LINESIZE - 1] = '\0';
	strncpy(l.ll_host, host, UT_HOSTSIZE);
	l.ll_host[UT_HOSTSIZE - 1] = '\0';

	l.ll_time = (time_t)t;
	// printf("l: ll_line '%s', ll_host '%s', ll_time %d\n", l.ll_line, l.ll_host, l.ll_time);

	/* Write lastlog entry at fixed offset (uid * sizeof(struct lastlog) */
	if (NULL != (f = fopen("/var/log/lastlog", "rw+"))) {
		if (!fseek(f, (uid * sizeof(struct lastlog)), SEEK_SET)) {
			int fd = fileno(f);
			if (write(fd, &l, sizeof(l)) == sizeof(l)) {
				retval = 1;
				// int32_t stat = system("echo ---- lastlog ----; lastlog");
			}
		}
		fclose(f);
	}
	return retval;
}
