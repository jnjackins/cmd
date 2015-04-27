package main

// from github.com/docker/libcontainer/devices/number.go
func devNums(dev uint64) (major, minor uint32) {
	num := uint32(dev)
	major = (num >> 8) & 0xfff
	minor = (num & 0xff) | ((num >> 12) & 0xfff00)
	return
}
