package mlog

import "os"

var direct *Ts = New(nil)

func SetConfig(config *Config) {
	direct = New(config)
}

func Fatal(fields H) *Ts {
	if LevelFatal > direct.config.Level {
		return direct
	}
	fields["level"] = "FATAL"
	direct.logWithLevel(fields, 2)
	os.Exit(1)
	return direct
}

func Error(fields H) *Ts {
	if LevelError > direct.config.Level {
		return direct
	}
	fields["level"] = "ERROR"
	return direct.logWithLevel(fields, 2)
}

func Warn(fields H) *Ts {
	if LevelWarn > direct.config.Level {
		return direct
	}
	fields["level"] = "WARN"
	return direct.logWithLevel(fields, 2)
}

func Info(fields H) *Ts {
	if LevelInfo > direct.config.Level {
		return direct
	}
	fields["level"] = "INFO"
	return direct.logWithLevel(fields, 2)
}

func Debug(fields H) *Ts {
	if LevelDebug > direct.config.Level {
		return direct
	}
	fields["level"] = "DEBUG"
	return direct.logWithLevel(fields, 2)
}

func Trace(fields H) *Ts {
	if LevelTrace > direct.config.Level {
		return direct
	}
	fields["level"] = "TRACE"
	return direct.logWithLevel(fields, 2)
}
