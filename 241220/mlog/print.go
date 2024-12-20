package mlog

import "os"

func (e *mlog) Fatal(fields H) *mlog {
	if LevelFatal > e.config.Level {
		return e
	}
	fields["level"] = "FATAL"
	e.logWithLevel(fields, 2)
	os.Exit(1)
	return e
}

func (e *mlog) Error(fields H) *mlog {
	if LevelError > e.config.Level {
		return e
	}
	fields["level"] = "ERROR"
	return e.logWithLevel(fields, 2)
}

func (e *mlog) Warn(fields H) *mlog {
	if LevelWarn > e.config.Level {
		return e
	}
	fields["level"] = "WARN"
	return e.logWithLevel(fields, 2)
}

func (e *mlog) Info(fields H) *mlog {
	if LevelInfo > e.config.Level {
		return e
	}
	fields["level"] = "INFO"
	return e.logWithLevel(fields, 2)
}

func (e *mlog) Debug(fields H) *mlog {
	if LevelDebug > e.config.Level {
		return e
	}
	fields["level"] = "DEBUG"
	return e.logWithLevel(fields, 2)
}

func (e *mlog) Trace(fields H) *mlog {
	if LevelTrace > e.config.Level {
		return e
	}
	fields["level"] = "TRACE"
	return e.logWithLevel(fields, 2)
}
