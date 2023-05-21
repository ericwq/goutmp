# goutmp

This is a modified golang module which support utmpx API. This is a temporary golang module. The next step is to create a pure golang client module to support [utmps](https://skarnet.org/software/utmps/). Currenly, The implementation is only a wrapper for `utmps` C client library.

## original goutmp README.md
[goutmp](https://gogs.blitter.com/RLabs/goutmp) - Minimal bindings to C stdlib pututmpx(), getutmpx() (/var/log/wtmp) and /var/log/lastlog

Any Go program which allows user shell access should update the standard UNIX files which track user sessions: /var/log/wtmp (for the 'w' and 'who' commands), and /var/log/lastlog (the 'last' and 'lastlog' commands).

```sh
$ go doc -all
package goutmp // import "blitter.com/go/goutmp"

Golang bindings for basic login/utmp accounting

FUNCTIONS

func GetHost(addr string) (h string)
    return remote client hostname or IP if host lookup fails addr is expected to
    be of the format given by net.Addr.String() eg., "127.0.0.1:80" or "[::1]:80"

func Put_lastlog_entry(app, usr, ptsname, host string)
    Put the login app, username and originating host/IP to lastlog

func Unput_utmp(entry UtmpEntry)
    Remove a username/host entry from utmp


TYPES

type UtmpEntry struct {
	// Has unexported fields.
}
    UtmpEntry wraps the C struct utmp

func Put_utmp(user, ptsName, host string) UtmpEntry
    Put a username and the originating host/IP to utmp
```


