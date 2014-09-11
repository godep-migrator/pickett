// Interface functions for the Global and Local loggers for each type.
// This was auto generated with makeInterface.py.

package logit

// Utility for FINEST log messages
// This behaves like Logf but with the FINEST log level.
func Finestf(format string, args ...interface{}) {
	Global.intLogf(FINEST, format, args...)
}

// Utility for FINEST log messages
// This behaves like Logln but with the FINEST log level.
func Finestln(args ...interface{}) {
	Global.intLogln(FINEST, args...)
}

// Utility for FINEST log messages
// This behaves like Logc but with the FINEST log level.
func Finestc(closure func() string) {
	Global.intLogc(FINEST, closure)
}

// Utility for FINE log messages
// This behaves like Logf but with the FINE log level.
func Finef(format string, args ...interface{}) {
	Global.intLogf(FINE, format, args...)
}

// Utility for FINE log messages
// This behaves like Logln but with the FINE log level.
func Fineln(args ...interface{}) {
	Global.intLogln(FINE, args...)
}

// Utility for FINE log messages
// This behaves like Logc but with the FINE log level.
func Finec(closure func() string) {
	Global.intLogc(FINE, closure)
}

// Utility for DEBUG log messages
// This behaves like Logf but with the DEBUG log level.
func Debugf(format string, args ...interface{}) {
	Global.intLogf(DEBUG, format, args...)
}

// Utility for DEBUG log messages
// This behaves like Logln but with the DEBUG log level.
func Debugln(args ...interface{}) {
	Global.intLogln(DEBUG, args...)
}

// Utility for DEBUG log messages
// This behaves like Logc but with the DEBUG log level.
func Debugc(closure func() string) {
	Global.intLogc(DEBUG, closure)
}

// Utility for TRACE log messages
// This behaves like Logf but with the TRACE log level.
func Tracef(format string, args ...interface{}) {
	Global.intLogf(TRACE, format, args...)
}

// Utility for TRACE log messages
// This behaves like Logln but with the TRACE log level.
func Traceln(args ...interface{}) {
	Global.intLogln(TRACE, args...)
}

// Utility for TRACE log messages
// This behaves like Logc but with the TRACE log level.
func Tracec(closure func() string) {
	Global.intLogc(TRACE, closure)
}

// Utility for INFO log messages
// This behaves like Logf but with the INFO log level.
func Infof(format string, args ...interface{}) {
	Global.intLogf(INFO, format, args...)
}

// Utility for INFO log messages
// This behaves like Logln but with the INFO log level.
func Infoln(args ...interface{}) {
	Global.intLogln(INFO, args...)
}

// Utility for INFO log messages
// This behaves like Logc but with the INFO log level.
func Infoc(closure func() string) {
	Global.intLogc(INFO, closure)
}

// Utility for WARNING log messages
// This behaves like Logf but with the WARNING log level.
func Warningf(format string, args ...interface{}) {
	Global.intLogf(WARNING, format, args...)
}

// Utility for WARNING log messages
// This behaves like Logln but with the WARNING log level.
func Warningln(args ...interface{}) {
	Global.intLogln(WARNING, args...)
}

// Utility for WARNING log messages
// This behaves like Logc but with the WARNING log level.
func Warningc(closure func() string) {
	Global.intLogc(WARNING, closure)
}

// Utility for ERROR log messages
// This behaves like Logf but with the ERROR log level.
func Errorf(format string, args ...interface{}) {
	Global.intLogf(ERROR, format, args...)
}

// Utility for ERROR log messages
// This behaves like Logln but with the ERROR log level.
func Errorln(args ...interface{}) {
	Global.intLogln(ERROR, args...)
}

// Utility for ERROR log messages
// This behaves like Logc but with the ERROR log level.
func Errorc(closure func() string) {
	Global.intLogc(ERROR, closure)
}

// Utility for CRITICAL log messages
// This behaves like Logf but with the CRITICAL log level.
func Criticalf(format string, args ...interface{}) {
	Global.intLogf(CRITICAL, format, args...)
}

// Utility for CRITICAL log messages
// This behaves like Logln but with the CRITICAL log level.
func Criticalln(args ...interface{}) {
	Global.intLogln(CRITICAL, args...)
}

// Utility for CRITICAL log messages
// This behaves like Logc but with the CRITICAL log level.
func Criticalc(closure func() string) {
	Global.intLogc(CRITICAL, closure)
}

// Utility for FINEST log messages
// This behaves like Logf but with the FINEST log level.
func (log *Logger) Finestf(format string, args ...interface{}) {
	log.intLogf(FINEST, format, args...)
}

// Utility for FINEST log messages
// This behaves like Logln but with the FINEST log level.
func (log *Logger) Finestln(args ...interface{}) {
	log.intLogln(FINEST, args...)
}

// Utility for FINEST log messages
// This behaves like Logc but with the FINEST log level.
func (log *Logger) Finestc(closure func() string) {
	log.intLogc(FINEST, closure)
}

// Utility for FINE log messages
// This behaves like Logf but with the FINE log level.
func (log *Logger) Finef(format string, args ...interface{}) {
	log.intLogf(FINE, format, args...)
}

// Utility for FINE log messages
// This behaves like Logln but with the FINE log level.
func (log *Logger) Fineln(args ...interface{}) {
	log.intLogln(FINE, args...)
}

// Utility for FINE log messages
// This behaves like Logc but with the FINE log level.
func (log *Logger) Finec(closure func() string) {
	log.intLogc(FINE, closure)
}

// Utility for DEBUG log messages
// This behaves like Logf but with the DEBUG log level.
func (log *Logger) Debugf(format string, args ...interface{}) {
	log.intLogf(DEBUG, format, args...)
}

// Utility for DEBUG log messages
// This behaves like Logln but with the DEBUG log level.
func (log *Logger) Debugln(args ...interface{}) {
	log.intLogln(DEBUG, args...)
}

// Utility for DEBUG log messages
// This behaves like Logc but with the DEBUG log level.
func (log *Logger) Debugc(closure func() string) {
	log.intLogc(DEBUG, closure)
}

// Utility for TRACE log messages
// This behaves like Logf but with the TRACE log level.
func (log *Logger) Tracef(format string, args ...interface{}) {
	log.intLogf(TRACE, format, args...)
}

// Utility for TRACE log messages
// This behaves like Logln but with the TRACE log level.
func (log *Logger) Traceln(args ...interface{}) {
	log.intLogln(TRACE, args...)
}

// Utility for TRACE log messages
// This behaves like Logc but with the TRACE log level.
func (log *Logger) Tracec(closure func() string) {
	log.intLogc(TRACE, closure)
}

// Utility for INFO log messages
// This behaves like Logf but with the INFO log level.
func (log *Logger) Infof(format string, args ...interface{}) {
	log.intLogf(INFO, format, args...)
}

// Utility for INFO log messages
// This behaves like Logln but with the INFO log level.
func (log *Logger) Infoln(args ...interface{}) {
	log.intLogln(INFO, args...)
}

// Utility for INFO log messages
// This behaves like Logc but with the INFO log level.
func (log *Logger) Infoc(closure func() string) {
	log.intLogc(INFO, closure)
}

// Utility for WARNING log messages
// This behaves like Logf but with the WARNING log level.
func (log *Logger) Warningf(format string, args ...interface{}) {
	log.intLogf(WARNING, format, args...)
}

// Utility for WARNING log messages
// This behaves like Logln but with the WARNING log level.
func (log *Logger) Warningln(args ...interface{}) {
	log.intLogln(WARNING, args...)
}

// Utility for WARNING log messages
// This behaves like Logc but with the WARNING log level.
func (log *Logger) Warningc(closure func() string) {
	log.intLogc(WARNING, closure)
}

// Utility for ERROR log messages
// This behaves like Logf but with the ERROR log level.
func (log *Logger) Errorf(format string, args ...interface{}) {
	log.intLogf(ERROR, format, args...)
}

// Utility for ERROR log messages
// This behaves like Logln but with the ERROR log level.
func (log *Logger) Errorln(args ...interface{}) {
	log.intLogln(ERROR, args...)
}

// Utility for ERROR log messages
// This behaves like Logc but with the ERROR log level.
func (log *Logger) Errorc(closure func() string) {
	log.intLogc(ERROR, closure)
}

// Utility for CRITICAL log messages
// This behaves like Logf but with the CRITICAL log level.
func (log *Logger) Criticalf(format string, args ...interface{}) {
	log.intLogf(CRITICAL, format, args...)
}

// Utility for CRITICAL log messages
// This behaves like Logln but with the CRITICAL log level.
func (log *Logger) Criticalln(args ...interface{}) {
	log.intLogln(CRITICAL, args...)
}

// Utility for CRITICAL log messages
// This behaves like Logc but with the CRITICAL log level.
func (log *Logger) Criticalc(closure func() string) {
	log.intLogc(CRITICAL, closure)
}
