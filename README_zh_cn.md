## 日志包的使用
日志的数量在正式过程中应不能太多，太大将可能影响程序性能;<br/>
在接口调用中，日志器采用了异步缓存的功能以提高日志的执行效率，缓存值由具体适配置器决定，默认为10条缓存。<br/>
日志在使用结束后应关闭，以确保缓存的数据能够输出，如果缓存的数据不能输出，该错误由适配器进行转存处理。<br/>

## 关于日志级别的调用
DEBUG <br/>
调试信息,仅调试时输出.<br/>
由平台控制,可控制台输出,相当于fmt.Print输出.<br/>

INFO <br/>
程序运行状态信息.<br/>
如启动、停止、重连等信息，体现了程序环境的变更状态。<br/>

WARN <br/>
程序异常信息。
本类别不影响程序继续使用,但其结果可能会引出潜在的重大问题。<br/>
例如：某些未知的错误；网络连接错误(但重新连接后可自动修复), 连接超时等行为。<br/>
此类异常在一段时间如果出现过多，那么应该分析其中的原因，例如可能存在的问题：<br/>
被攻击、硬件老化、硬件达到了承载上限、对方服务出现异常等问题。<br/>
日志系统建议发送一封邮件到相关人员。<br/>

ERROR<br/>
程序致命的错误信息。
此错误将影响到正常逻辑, 甚至平台因此而恐慌、停止服务的行为.<br/>
例如：数据库不可用、充值不可用、短信不可用、vos不可用等需要及时处理的行为都可定义为此类别。<br/>
日志系统建议发送一封邮件、实时短信或电话(或者其他实时联系方式)通知至相关人员<br/>


FATAL<br/>
程序非正常结束。<br/>
日志系统建议调用所有实时联系方式联系相关人员处理。<br/>

## 默认包使用例子:

```
import (
  "github.com/gwaylib/log"
)

func main() {
  log.Info("OK")
}
```

## 定制一个日志器
```
import (
  "github.com/gwaylib/log"
  "github.com/gwaylib/logger"
)

func init() {
  // make a custom logger 
  level, _ := strconv.Atoi(os.Getenv(log.GWAYLIB_LOG_LEVEL_NAME))
  adapter = []logger.Adapter{stdio.New(os.Stdout)}
  log.Log = logger.New(&logger.DefaultContext, "appname", proto.Levev(level), adapter...)
}

func main() {
  log.Info("OK")
}
```
