package mlog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type H map[string]interface{}

type Ts struct {
	config *Config

	H      map[string]any
	logger *log.Logger
}

type Config struct {
	Stdout      bool     `json:"stdout" yaml:"stdout"`
	Level       int      `json:"level" yaml:"level"`
	FileName    string   `json:"filename" yaml:"filename"`
	OrderedKeys []string `json:"orderedkeys" yaml:"orderedkeys"`
	CallerClip  string   `json:"callerclip" yaml:"callerclip"`

	// lumberjack.Logger fields
	RotateMaxSize    int  `json:"rotatemaxsize" yaml:"rotatemaxsize"`
	RotateMaxAge     int  `json:"rotatemaxage" yaml:"rotatemaxage"`
	RotateMaxBackups int  `json:"rotatemaxbackups" yaml:"rotatemaxbackups"`
	RotateLocalTime  bool `json:"rotatelocaltime" yaml:"rotatelocaltime"`
	RotateCompress   bool `json:"rotatecompress" yaml:"rotatecompress"`
}

func New(config *Config) *Ts {
	if config == nil {
		config = NewConfig()
		config.FileName = ""
	}

	ref := &Ts{config: config}
	if ref.config.FileName != "" {
		logger := &lumberjack.Logger{
			Filename:   ref.config.FileName,
			MaxSize:    ref.config.RotateMaxSize,
			MaxBackups: ref.config.RotateMaxBackups,
			MaxAge:     ref.config.RotateMaxAge,
			Compress:   ref.config.RotateCompress,
			LocalTime:  ref.config.RotateLocalTime,
		}
		ref.logger = log.New(logger, "", 0)
	}
	return ref
}

func NewConfig() *Config {
	execName := filepath.Base(os.Args[0])
	fileName := fmt.Sprintf("/var/log/%s.log", execName)
	return &Config{
		Stdout:           true,
		Level:            3,
		FileName:         fileName,
		OrderedKeys:      []string{"time", "level", "msg", "info", "error", "warn", "data", "flags"},
		CallerClip:       "",
		RotateMaxSize:    10,
		RotateMaxAge:     10,
		RotateMaxBackups: 10,
		RotateLocalTime:  true,
		RotateCompress:   true,
	}
}

func (e *Ts) logWithLevel(fields H, callDepth int) *Ts {
	e.setCaller(fields, callDepth+1)

	coloredOutput := colorizeJSONValues(fields, e.config.OrderedKeys)

	if e.config.FileName != "" {
		e.logger.Println(coloredOutput)
	}
	if e.config.Stdout {
		fmt.Println(coloredOutput)
	}

	return e
}

func colorizeJSONValues(fields H, orderedKeys []string) string {
	resetColor := "\033[0m"

	levelColors := map[string]string{
		"FATAL": "\033[95m",
		"ERROR": "\033[91m",
		"WARN":  "\033[93m",
		"INFO":  "\033[92m",
		"DEBUG": "\033[94m",
		"TRACE": "\033[90m",
	}

	keyColors := map[string]string{
		"time":  "\033[30m",
		"msg":   "\033[34m",
		"error": "\033[31m",
		"warn":  "\033[33m",
		"info":  "\033[34m",
		"data":  "\033[32m",
		"other": "\033[36m",
		"call":  "\033[35m",
	}

	var otherKeys []string
	for key := range fields {
		if key != "call" && !slices.Contains(orderedKeys, key) {
			otherKeys = append(otherKeys, key)
		}
	}
	sort.Strings(otherKeys)

	allKeys := append(orderedKeys, otherKeys...)

	if _, exists := fields["call"]; exists {
		allKeys = append(allKeys, "call")
	}

	var builder strings.Builder
	builder.WriteByte('{')

	first := true
	for _, key := range allKeys {
		value, exists := fields[key]
		if !exists {
			continue
		}

		if !first {
			builder.WriteByte(',')
		}
		first = false

		builder.WriteByte('"')
		builder.WriteString(key)
		builder.WriteString(`":`)

		valueBytes, err := json.Marshal(value)
		if err != nil {
			valueBytes = []byte(`"` + fmt.Sprintf("%v", value) + `"`)
		}

		if key == "level" {
			level := fmt.Sprintf("%v", value)
			if color, ok := levelColors[level]; ok {
				builder.WriteString(color)
				builder.Write(valueBytes)
				builder.WriteString(resetColor)
				continue
			}
		}

		if color, ok := keyColors[key]; ok {
			builder.WriteString(color)
			builder.Write(valueBytes)
			builder.WriteString(resetColor)
		} else {
			builder.WriteString(keyColors["other"])
			builder.Write(valueBytes)
			builder.WriteString(resetColor)
		}
	}
	builder.WriteByte('}')

	return builder.String()
}
