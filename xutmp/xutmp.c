#include "xutmp.h"
#include <lastlog.h>
#include <stdlib.h>

// typedef char char_t;
struct utmpx* res = NULL;

// #define DEV_PREFIX "/dev/"
// #define DEV_PREFIX_LEN (sizeof(DEV_PREFIX) - 1)

#define MIN(a_, b_) (((a_) < (b_)) ? (a_) : (b_))

static int write_uwtmp_record(const char* user,
							  const char* termName,
							  const char* host,
							  pid_t pid,
							  int add) {
	struct utmpx ut;
	struct timeval tv;
	size_t len, offset;

	memset(&ut, 0, sizeof(ut));

	memset(&tv, 0, sizeof(tv));
	(void)gettimeofday(&tv, 0);

	if (user) {  // Logout record has null username
		len = strlen(user);
		memcpy(ut.ut_user, user, MIN(sizeof(ut.ut_user), len));
	}

	if (host) {  // Logout record has null host
		len = strlen(host);
		memcpy(ut.ut_host, host, MIN(sizeof(ut.ut_host), len));
	}

	// len = strlen(term);
	// memcpy(ut.ut_line, term, MIN(sizeof(ut.ut_line), len));
	//
	// offset = len <= sizeof(ut.ut_id) ? 0 : len - sizeof(ut.ut_id);
	// memcpy(ut.ut_id, term + offset, len - offset);

	// Set ut_line and ut_id based on the terminal associated with 'stdin'. This
	// code assumes terminals named "/dev/[pt]t[sy]*". The "/dev/" dirname is 5
	// characters; the "[pt]t[sy]" filename prefix is 3 characters (making 8
	// characters in all).

	len = strlen(termName + 5);
	memcpy(ut.ut_line, termName + 5, MIN(sizeof(ut.ut_line), len));

	len = strlen(termName + 8);
	memcpy(ut.ut_id, termName + 8, MIN(sizeof(ut.ut_id), len));

	if (add)
		ut.ut_type = USER_PROCESS;
	else
		ut.ut_type = DEAD_PROCESS;

	ut.ut_pid = pid;

	ut.ut_tv.tv_sec = (__typeof__(ut.ut_tv.tv_sec))tv.tv_sec;
	ut.ut_tv.tv_usec = (__typeof__(ut.ut_tv.tv_usec))tv.tv_usec;

	setutxent();
	if (!pututxline(&ut))
		return EXIT_FAILURE;
	// fatal_error("pututline: %s", strerror(errno));
	endutxent();

	(void)updwtmp(_PATH_WTMP, &ut);

	// debug_msg("utmp/wtmp record %s for terminal '%s'",
	// 	  add ? "added" : "removed", term);
	return EXIT_SUCCESS;
}

void pututmp(struct utmpx* ut, char* uname, char* ptsname, char* host) {
	// printf("effective GID=%u\n", getegid());
	// system("echo ---- pre ----;who");
	memset(ut, 0, sizeof(struct utmpx));

	ut->ut_type = USER_PROCESS;  // This is a user login
	strncpy(ut->ut_user, uname, sizeof(ut->ut_user));
	time((time_t*)&(ut->ut_tv.tv_sec));  // Stamp with current time
	ut->ut_pid = getpid();

	// Set ut_line and ut_id based on the terminal associated with 'stdin'. This code assumes
	// terminals named "/dev/[pt]t[sy]*". The "/dev/" dirname is 5 characters; the "[pt]t[sy]"
	// filename prefix is 3 characters (making 8 characters in all).

	// devName = ttyname(STDIN_FILENO);
	// if (devName == NULL)
	// 	errExit("ttyname");
	// if (strlen(devName) <= 8) // Should never happen
	// 	fatal("Terminal name is too short: %s", devName);
	strncpy(ut->ut_line, ptsname + 5, sizeof(ut->ut_line));
	strncpy(ut->ut_id, ptsname + 8, sizeof(ut->ut_id));
	strcpy(ut->ut_host, host);

	setutxent();               // Rewind to start of utmp file
	pututxline(ut);            // Overwrite previous utmp record
	updwtmpx(_PATH_WTMP, ut);  // Append login record to wtmp
	endutxent();
	// system("echo ---- post ----;who");
}

void unpututmp(struct utmpx* ut) {
	ut->ut_type = DEAD_PROCESS;              // Required for logout record
	time((time_t*)&(ut->ut_tv.tv_sec));      // Stamp with logout time
	memset(&(ut->ut_user), 0, UT_NAMESIZE);  // Logout record has null username
	setutxent();
	pututxline(ut);
	updwtmpx(_PATH_WTMP, ut);  // Append logout record to wtmp
	endutxent();

	// system("echo ---- cleanup ----;who; last");
}

struct utmpx* getutmp() {
	if (res != NULL)  // If 'res' was set via a previous call
		memset(res, 0, sizeof(struct utmpx));
	res = getutxent();
	if (res == NULL) {
		return NULL;
	}

	// unsigned char* charPtr = (unsigned char*)res;
	// int i;
	// int start = 32+4+4+2;
	// int end = start+4;  // sizeof(struct utmpx); )
	// for (i = start; i < end; i++)
	// 	printf("%02x ", charPtr[i]);
	// printf("\n");
	//
	// printf(
	// 	"[ C] type=%d; pid=%d; line=%s, id=%.4s; user=%s; host=%s; exit={%u %u}; session=%d; "
	// 	"time={%ld %ld}\n",
	// 	res->ut_type, res->ut_pid, res->ut_line, res->ut_id, res->ut_user, res->ut_host,
	// 	res->ut_exit.e_termination, res->ut_exit.e_exit, res->ut_session, res->ut_tv.tv_sec,
	// 	res->ut_tv.tv_usec);
	return res;
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
