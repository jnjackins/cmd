package main

import "strings"

type record struct {
	s      string
	fields []string
}

func splitRecord(line string) record {
	return record{
		s:      line,
		fields: strings.Fields(line),
	}
}

func (rec record) String() string {
	return rec.s
}
