# Unix userspace programs in Go

This project aims to provide a full, standalone Unix userspace environment implemented in Go. Just add a kernel. Compatible with Linux, and (hopefully) other Unix kernels such as FreeBSD.

## Why?
Just for fun. Or see [this paper](http://harmful.cat-v.org/cat-v/unix_prog_design.pdf). Or because it is easier to write correct and understandable programs in Go than in C. You decide!

## Contents
The included programs are functional, and should be sufficient to boot a working OS with the only other requirements being a kernel, bootloader, and an appropriately laid out filesystem. The most notable missing piece is a usable shell interpreter.

Most programs in this repository are minimalist versions of the usual ones, and may not have all the options that you expect. Those that differ significantly in their usage have been given different names. For example:

* _hget_ provides similar functionality to curl(1) and wget(1) 
* _ipconfig_ combines functionality from both ifconfig(8) and route(8)
* _sub_ provides substitution similar to s/a/b/ commands in sed(1), but using [re2 syntax](https://golang.org/s/re2syntax). Referencing parenthesized submatches in the replacement string with \1, \2, etc. is supported.

## TODO
* Provide a VM image
* Provide a decent shell interpreter
* Write a short README for each directory
