package mlog

import "os"

func (e *Ts) Fatal(fields H) *Ts {
	if LevelFatal > e.config.Level {
		return e
	}
	fields["level"] = "FATAL"
	e.logWithLevel(fields, 2)
	os.Exit(1)
	return e
}

func (e *Ts) Error(fields H) *Ts {
	if LevelError > e.config.Level {
		return e
	}
	fields["level"] = "ERROR"
	return e.logWithLevel(fields, 2)
}

func (e *Ts) Warn(fields H) *Ts {
	if LevelWarn > e.config.Level {
		return e
	}
	fields["level"] = "WARN"
	return e.logWithLevel(fields, 2)
}

func (e *Ts) Info(fields H) *Ts {
	if LevelInfo > e.config.Level {
		return e
	}
	fields["level"] = "INFO"
	return e.logWithLevel(fields, 2)
}

func (e *Ts) Debug(fields H) *Ts {
	if LevelDebug > e.config.Level {
		return e
	}
	fields["level"] = "DEBUG"
	return e.logWithLevel(fields, 2)
}

func (e *Ts) Trace(fields H) *Ts {
	if LevelTrace > e.config.Level {
		return e
	}
	fields["level"] = "TRACE"
	return e.logWithLevel(fields, 2)
}
