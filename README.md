# goutmp

This is a modified golang module which support utmpx API. This is a temporary golang module. The next step is to create a pure golang client module to support [utmps](https://skarnet.org/software/utmps/). Currenly, The implementation is only a wrapper for `utmps` C client library.

## inline C or stand alone C module
We use the following cgo derective and inline C functions to implement the wrapper. Inline C functions is more easy to build than stand alone C module.
```c
// #cgo pkg-config: utmps skalibs
```

Please use the following commands to build the stand alone C module, either staticly or dyanamicly.
```sh 
$ cd ./xutmp/
$ gcc -I/usr/include/utmps -lutmps -lskarnet -c -o xutmp.o xutmp.c
$ ar rcs libxutmp.a xutmp.o
```

```sh
$ gcc -shared -I/usr/include/utmps -lutmps -lskarnet -o libxutmp.so xutmp.c
```

The following cgo derective is for stand alone C module.
```c
/*
#cgo CFLAGS: -I./xutmp
#cgo LDFLAGS: -L${SRCDIR}/xutmp -lxutmp

#include "xutmp.h"
*/
```

## prepare the utmps environment 

Refer to [alpine container with openrc support](https://github.com/ericwq/s6), run the following command to start the container.
```sh
% docker run --env TZ=Asia/Shanghai --tty --privileged --volume /sys/fs/cgroup:/sys/fs/cgroup:ro \
  -h openrc-ssh --name openrc-ssh -d -p 5022:22 openrc-ssh:0.1.0
```

Install `go` SDK and `utmps-dev` package. Note, this command require root privilege.
```sh
% apk add go utmps-dev
```

Run `setup-utmp` script to prepare the container. Note, this command require root privilege.
```shared
% setup-utmp
```

Restart the container. Run the following command to make sure everything works.
```
openrc-ssh:~# pstree -p
init(1)-+-s6-ipcserverd(154)
        |-s6-ipcserverd(217)
        |-s6-ipcserverd(245)
        `-sshd(190)---sshd(286)---ash(288)---pstree(292)
openrc-ssh:~# who
root            pts/1           00:00   May 29 22:48:09  172.17.0.1
openrc-ssh:~# last
USER       TTY            HOST               LOGIN        TIME
root       pts/1          172.17.0.1         May 29 22:48
```

Now, you can build your application and run it according to the following section.

## how to set effective GID for your service
You has a service and want that service has the privileges to access `utmps` service. Then you need to set the effective GID for your service to be `utmp`. The `utmps` service require effective GID of `utmp`. Refer to [The utmps-utmpd program](https://skarnet.org/software/utmps/utmps-utmpd.html) for detail.

Let's say your service program is `prog2`. You need set GID for `prog2`. Let's assume that user `ide` belongs to two groups: `develop` and `utmp`. You can use `$ adduser ide utmp` command to achive it.

- first, change the group of `prog2` to `utmp`.
- second, set-GID for `prog2`.
- finally, if you run the `prog2` program, it's effective GID will be `utmp`.

```sh
openrc-ssh:/tmp$ ls -al
total 2096
drwxrwxrwt    1 root     root          4096 May 27 20:38 .
drwxr-xr-x    1 root     root          4096 May 27 18:12 ..
-rwxr-xr-x    1 ide      develop    2137512 May 27 20:38 prog2
openrc-ssh:/tmp$ chgrp utmp prog2
openrc-ssh:/tmp$ ls -al
total 2096
drwxrwxrwt    1 root     root          4096 May 27 20:38 .
drwxr-xr-x    1 root     root          4096 May 27 18:12 ..
-rwxr-xr-x    1 ide      utmp       2137512 May 27 20:38 prog2
openrc-ssh:/tmp$ chmod g+s prog2
openrc-ssh:/tmp$ ls -al
total 2096
drwxrwxrwt    1 root     root          4096 May 27 20:38 .
drwxr-xr-x    1 root     root          4096 May 27 18:12 ..
-rwxr-sr-x    1 ide      utmp       2137512 May 27 20:38 prog2
```

Please refer to [s6-setuidgid](https://skarnet.org/software/s6/s6-setuidgid.html) to accomplish the above work in a single command. Note: in docker environment, mounted local file system does not support set UID/GID operation.

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


