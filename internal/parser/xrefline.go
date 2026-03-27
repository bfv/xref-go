package parser

import (
	"strconv"
	"strings"
)

// XrefLine represents a single parsed line from an xref file.
type XrefLine struct {
	CompiledFile   string
	ContainingFile string
	LineNumber     int
	Type           string
	Info           string
}

// ParseLine parses a raw xref line into the XrefLine fields.
func (xl *XrefLine) ParseLine(line string) {
	entries := strings.Fields(strings.TrimSpace(line))
	if len(entries) < 4 {
		return
	}

	xl.CompiledFile = entries[0]
	xl.ContainingFile = entries[1]
	xl.LineNumber, _ = strconv.Atoi(entries[2])
	xl.Type = entries[3]

	if len(entries) > 4 {
		xl.Info = strings.Join(entries[4:], " ")
	}
}
