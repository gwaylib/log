## Using gwaylib/log
The number of logs should not be too many in the release program, which may affect the program performance.<br/>
In the call, the logger adopts the function of asynchronous cache to improve the execution efficiency of the log. The cache value is determined by the adapter implements, and the default is 10 items cache.<br/>
The log should be closed after use to ensure that the cached data can be output. If the cached data cannot be output, the error will be transferred by the adapter.<br/>

## About using log level
DEBUG <br/>
Debugging information, output only during debugging. <br/>
It is controlled by the platform and can be output to the console, which is equivalent to fmt.Print output.<br/>

INFO <br/>
Program running status change information.<br/>
Such as start, stop, reconnection and other information, which reflects the change state of the program environment.<br/>

WARN <br/>
Program exception information.<br/>
This category does not affect the continued use of the program, but its results may lead to potential major problems.<br/>
For example: some unknown errors; Network connection error (but it can be repaired automatically after reconnection), connection timeout, etc.<br/>
If there are too many such exceptions in a period of time, the causes should be analyzed, such as possible problems:<br/>
Being attacked, hardware aging, hardware reaching the upper limit of bearing capacity, abnormal service of the other party, etc. <br/>
The log adapter recommends sending an email to relevant personnel.<br/>

ERROR<br/>
Program fatal error message. <br/>
This error will affect the normal logic, and even the platform panic and stop the service. <br/>
For example, behaviors that need to be handled in time, such as database unavailable, recharge unavailable, SMS unavailable, and Vos unavailable, can be defined as this category. <br/>
The log adapter recommends sending an email, real-time SMS or telephone (or other real-time contact information) to notify relevant personnel.<br/>

FATAL<br/>
The program ended abnormally.<br/>
The log adapter recommends calling all real-time contact information and contacting relevant personnel for processing.<br/>

## Default use case:

```
import (
  "github.com/gwaylib/log"
)

func main() {
  log.Info("OK")
}
```

## Custom a logger
```
import (
  "github.com/gwaylib/log"
  "github.com/gwaylib/logger"
)

var lg *logger.Log

func init() {
  // make a custom logger 
  level, _ := strconv.Atoi(os.Getenv(log.GWAYLIB_LOG_LEVEL))
  adapter := []logger.Adapter{stdio.New(os.Stdout)}
  callerDepth := 4 // if callerDepth == 0, close the file path caller

  // replace the default logger
  log.Log = logger.New(&logger.DefaultContext, "appname", callerDepth, proto.Level(level), adapter...)

  // or
  // make a new package log
  lg = logger.New(&logger.DefaultContext, "appname", callerDepth, proto.Level(level), adapter...)

  // or
  // make a new package default log
  //lg = logger.NewDefaultLogger(loggerName, adapter...)

  // or 
  // copy the log.go file in local and fix it
}

func main() {
  log.Info("OK")
  lg.Info("OK")
}
```
