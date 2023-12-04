package dev

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	
	"github.com/creamsensation/cp/internal/constant/pkg"
)

type ErrorTrace struct {
	Line int
	Rows []string
	Path string
}

const (
	errorTraceRange = 3
)

func GetErrorTrace() []ErrorTrace {
	stackSlice := make([]byte, 1024)
	s := runtime.Stack(stackSlice, false)
	result := make([]ErrorTrace, 0)
	wd, err := os.Getwd()
	if err != nil {
		return result
	}
	for _, item := range strings.Split(fmt.Sprintf("%s", stackSlice[0:s]), "\n") {
		if !strings.Contains(item, ".go") ||
			strings.Contains(item, "/internal/") ||
			strings.Contains(item, "/runtime/panic") ||
			strings.Contains(item, "/http/server") ||
			strings.Contains(
				item, fmt.Sprintf("/%s/", pkg.Name),
			) {
			continue
		}
		if strings.Contains(item, "/reflect/") {
			break
		}
		item = strings.TrimSpace(item)
		if strings.Contains(item, " ") {
			item = item[:strings.Index(item, " ")]
		}
		path := item[:strings.Index(item, ":")]
		line, err := strconv.Atoi(item[strings.Index(item, ":")+1:])
		if err != nil {
			continue
		}
		fileBts, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		rows := getFileContentSnippet(line, string(fileBts))
		result = append(
			result, ErrorTrace{
				Line: line,
				Rows: rows,
				Path: strings.TrimPrefix(path, wd),
			},
		)
	}
	return result
}

func getFileContentSnippet(line int, content string) []string {
	result := make([]string, 0)
	parts := strings.Split(content, "\n")
	partsLen := len(parts)
	for i, row := range parts {
		bottomLimit := line - errorTraceRange
		topLimit := line + errorTraceRange
		if bottomLimit < 0 {
			bottomLimit = 0
		}
		if topLimit > partsLen-1 {
			topLimit = partsLen - 1
		}
		if i > bottomLimit && i < topLimit {
			result = append(result, row)
		}
	}
	return result
}
