# Init

Init is a simple implementation Research Unix style [init(8)](http://www.unix.com/man-page/v7/8/init/), loosely based on the Research
Unix Version 7 init.

Notably, this implementation of init does not launch [getty(8)](http://www.unix.com/man-page/v7/8/getty/) processes, but rather uses per-tty
goroutines to launch [login(1)](http://www.unix.com/man-page/v7/1/login/) processes directly.

## Phases

Init loops through a series of phases, in the following order:

### 1. Shutdown 
The shutdown phase terminates any running session goroutines and kills all child processes.
In case of a panic, the behaviour of init is to recover and start execution again from Phase 1.
Thus, the shutdown phase is first to ensure that any running sessions are closed and all child processes are cleaned up.

### 2. Single-user
Launch a single [sh(1)](http://www.unix.com/man-page/v7/1/sh/) session with no authentication. To proceed to the next phase, send EOF
by pressing ctl-d.

### 3. Runcom (Run Commands)
Execute the [sh(1)](http://www.unix.com/man-page/v7/1/sh/) program /etc/rc. This can be used for system initialization, starting daemons,
etc.

### 4. Multi-user
Start a session per TTY. Switch between sessions by pressing ctl+alt+F1..FN. Each session
begins by launching [login(1)](http://www.unix.com/man-page/v7/1/login/). To end multi-user mode, close all sessions, and proceed to
the shutdown phase; send SIGHUP to init: ```kill -HUP 1```
