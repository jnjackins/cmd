package main

const (
	minorbits = 20
	minormask = (1 << minorbits) - 1
)

func devNums(dev int32) (major, minor uint32) {
	major = uint32(dev) >> minorbits
	minor = uint32(dev) & minormask
	return
}
