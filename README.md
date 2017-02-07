# Just Add a Kernel

This project aims to provide a full, standalone Unix userspace environment implemented in Go.  compatible with the Linux kernel, and hopefully other Unix-like kernels such as FreeBSD.

## Why?
Just for fun. Or see [this paper](http://harmful.cat-v.org/cat-v/unix_prog_design.pdf). Or because it is easier to write correct and understandable programs in Go than in C. You decide!

## Contents
The included programs are functional, and should be sufficient to boot a working OS with the only other requirements being a kernel, bootloader, and an appropriately laid out filesystem. The most notable missing piece is a usable shell interpreter.

Most programs in this repository are minimalist versions of the usual ones, and may not have all the options that you expect. Those that differ significantly in their usage have been given different names. For example:

* _hget_ provides similar functionality to curl(1) and wget(1) 
* _ipconfig_ combines functionality from both ifconfig(8) and route(8)

## TODO
* Provide a VM image
* Provide a decent shell interpreter
* Write a short README for each directory
