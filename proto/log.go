package proto

type Adapter interface {
	Put(log *Proto)
	Close()
}

// TODO: output the failed log to syslog
func FailLog(err error) {
	println(err.Error())
}

type Logger interface {
	SetOutputLevel(level Level)
	GetOutputLevel() Level

	SetAdapter(adapter []Adapter)
	GetAdapter() []Adapter

	Debug(msg ...interface{})
	Debugf(f string, msg ...interface{})

	Info(msg ...interface{})
	Infof(f string, msg ...interface{})

	Warn(msg ...interface{})
	Warnf(f string, msg ...interface{})

	Error(msg ...interface{})
	Errorf(f string, msg ...interface{})

	Fatal(msg ...interface{})
	Fatalf(f string, msg ...interface{})

	Exit(code int, msg ...interface{}) // exit with info level message
	ExitWithLevel(code int, level Level, msg ...interface{})

	// When exit program, need to call Close function to fush the log output.
	Close()

	// Compatibility interfaces
	// Same as Info level
	Print(msg ...interface{})
	Printf(f string, msg ...interface{})

	// Compatibility interfaces
	// Same as Fatal level
	Panic(msg ...interface{})
	Panicf(f string, msg ...interface{})
}
