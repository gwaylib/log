Example:

Create a log package in in your project.
```shell
mkdir $PRJ_ROOT/src/log
cp log.go $PRJ_ROOT/src/log
```

Change the hard code in your project
```golang 
var (
	level   = proto.LevelDebug // set the default log level
	adapter []logger.Adapter 
	Log     = New("your program name") // set your logger name
)

func init() {
	// implement the adapter what you need.
	adapter = []logger.Adapter{stdio.New(os.Stdout, os.Stderr)}
}
```

// Call your log package with default logger name.
```golang 
import "gomod-path/log" // replace gomod-path to your real package path.

func main(){
    log.Info("Hello")
}
```

// Call your log package with a defined logger name.
```golang 
import l "gomod-path/log" // replace gomod-path to your real package path.

var log = l.New("testing") // define a special logger name.

func main(){
    log.Info("Hello")
}
```

