package mlog

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

func (e *Ts) setCaller(fields H, callDepth int) {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		file = "unknown"
		line = 0
	}

	fields["call"] = e.pathClipping(fmt.Sprintf("%s:%d", file, line))

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fields["time"] = formattedTime
}

func (e *Ts) pathClipping(path string) string {
	if e.config.CallerClip == "" {
		if len(path) > 0 && path[0] == '/' {
			parts := strings.Split(path, "/")
			startIndex := len(parts) - 3
			if startIndex < 0 {
				startIndex = 0
			}
			path = "/" + strings.Join(parts[startIndex:], "/")
		}
		return path
	}
	return strings.Replace(path, e.config.CallerClip, "", -1)

}
