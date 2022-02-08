package fxrus

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/fx/fxevent"
)

// Logger is a wrapper for a logrus entry, that can be used as an fx event logger
type Logger struct {
	log logrus.FieldLogger
}

// Logger returns a new logger
func NewLogger(l logrus.FieldLogger) func() fxevent.Logger {
	return func() fxevent.Logger { return &Logger{log: l} }
}

// LogEvent logs an event
func (l *Logger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.log.WithFields(logrus.Fields{
			"hook":     "OnStart",
			"caller":   e.CallerName,
			"function": e.FunctionName,
		}).Debug("Executing hook")
	case *fxevent.OnStartExecuted:
		entry := l.log.WithFields(logrus.Fields{
			"hook":     "OnStart",
			"caller":   e.CallerName,
			"function": e.FunctionName,
			"runtime":  e.Runtime,
		})
		if e.Err != nil {
			entry.WithError(e.Err).Error("Failed to execute hook")
		} else {
			entry.Debug("Executed hook")
		}
	case *fxevent.OnStopExecuting:
		l.log.WithFields(logrus.Fields{
			"hook":     "OnStop",
			"caller":   e.CallerName,
			"function": e.FunctionName,
		}).Debug("Executing hook")
	case *fxevent.OnStopExecuted:
		entry := l.log.WithFields(logrus.Fields{
			"hook":     "OnStop",
			"caller":   e.CallerName,
			"function": e.FunctionName,
			"runtime":  e.Runtime,
		})
		if e.Err != nil {
			entry.WithError(e.Err).Error("Failed to execute hook")
		} else {
			entry.Debug("Executed hook")
		}
	case *fxevent.Supplied:
		entry := l.log.WithFields(logrus.Fields{
			"type": e.TypeName,
		})
		if e.Err != nil {
			entry.WithError(e.Err).Error("Failed to supply dependency")
		} else {
			entry.Debug("supplied dependency")
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.log.WithFields(logrus.Fields{"type": rtype}).Debug("Provided dependency")
		}
		if e.Err != nil {
			l.log.WithError(e.Err).Error("Error after options were applied")
		}
	case *fxevent.Invoking:
		l.log.WithFields(logrus.Fields{
			"function": e.FunctionName,
		}).Debug("Invoking function")
	case *fxevent.Invoked:
		entry := l.log.WithFields(logrus.Fields{
			"function": e.FunctionName,
			"trace":    e.Trace,
		})
		if e.Err != nil {
			entry.WithError(e.Err).Error("Failed to invoke function")
		} else {
			entry.Debug("Invoked function")
		}
	case *fxevent.Stopping:
		l.log.WithFields(logrus.Fields{
			"signal": e.Signal.String(),
		}).Debug("Stopping")
	case *fxevent.Stopped:
		if e.Err != nil {
			l.log.WithError(e.Err).Error("Failed to stop")
		} else {
			l.log.Info("Stopped")
		}
	case *fxevent.RollingBack:
		l.log.WithError(e.StartErr).Error("Rolling back")
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.log.WithError(e.Err).Error("Failed to roll back")
		} else {
			l.log.Debug("Rolled back")
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.log.WithError(e.Err).Error("Failed to start")
		} else {
			l.log.Info("Started")
		}
	case *fxevent.LoggerInitialized:
		entry := l.log.WithFields(logrus.Fields{
			"constructor": e.ConstructorName,
		})
		if e.Err != nil {
			entry.WithError(e.Err).Error("Failed to initialize logger")
		} else {
			entry.Debug("Initialized logger")
		}
	}
}
