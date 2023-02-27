package go_log

import (
	"io"
	"time"
)

type LogLevel string //日志级别

type TimeLayOut string //时间格式

type Color string //颜色

const (
	LoglevelTrace LogLevel = "TRACE"
	LoglevelDebug LogLevel = "DEBUG"
	LoglevelInfo  LogLevel = "INFO"
	LoglevelWarn  LogLevel = "WARN"
	LoglevelError LogLevel = "ERROR"
)

const (
	Reset             = "\033[0m"
	Red         Color = "\033[31m"
	Green       Color = "\033[32m"
	Yellow      Color = "\033[33m"
	Blue        Color = "\033[34m"
	Magenta     Color = "\033[35m"
	Cyan        Color = "\033[36m"
	White       Color = "\033[37m"
	BlueBold    Color = "\033[34;1m"
	MagentaBold Color = "\033[35;1m"
	RedBold     Color = "\033[31;1m"
	YellowBold  Color = "\033[33;1m"
)

const (

	// DefaultLayout 默认日志时间格式
	// 2006为Golang诞生时间，15是下午3点。帮助记忆的方法：1月2日3点4分5秒，2006年，-7时区，正好是1234567
	DefaultLayout   TimeLayOut = "2006-01-02 15:04:05.000"
	DateLayout      TimeLayOut = "2006-01-02"
	TimeLayout      TimeLayOut = "15:04:05"
	DateTimeLayout1 TimeLayOut = "2006-01-02-15-04-05"
	DateTimeLayout2 TimeLayOut = "20060102150405"
	DateTimeLayout3 TimeLayOut = "20060102_150405"
	DateTimeLayout4 TimeLayOut = "200601021504"
)

// LevelNum
//
//	@Description: 获取日志级别字符串
//	@receiver l
//	@return string
func (l LogLevel) LevelNum() int8 {
	switch l {
	case LoglevelTrace:
		return 1
	case LoglevelDebug:
		return 2
	case LoglevelInfo:
		return 3
	case LoglevelWarn:
		return 4
	case LoglevelError:
		return 5
	default:
		return -1
	}
}

// String
//
//	@Description: 获取时间字符串
//	@receiver t
//	@param layOut 时间格式，传入多个只使用第一个，不传使用默认格式 @See defaultLayout
//	@return string
func (t TimeLayOut) String(layOut ...TimeLayOut) string {
	if len(layOut) == 0 {
		return time.Now().Format(string(DefaultLayout))
	}
	return time.Now().Format(string(layOut[0]))

}

// WithColor
//
//	@Description: val及以后的值将以指定颜色开始，返回值后面再接值的时候也会带颜色
//	@receiver c
//	@param val
//	@return string
func (c Color) WithColor(val string) string {
	return string(c) + val
}

// WithColorEnd
//
//	@Description: val的值将会带颜色，不会影响后面的值
//	@receiver c
//	@param val
//	@return string
func (c Color) WithColorEnd(val string) string {
	return string(c) + val + Reset
}

// ILogger
// @Description: 日志的抽象接口
type ILogger interface {
	Trace(format string, msg ...any)
	// Debug Debug级别日志
	Debug(format string, msg ...any)
	// Info Info级别日志
	Info(format string, msg ...any)
	// Warn Warn级别日志
	Warn(format string, msg ...any)
	// Error Error级别日志
	Error(format string, msg ...any)
	// SetLogLevel 设置日志级别
	SetLogLevel(loglevel LogLevel)
	// SetLohWriter 设置输出流
	SetLohWriter(writer io.Writer)
	// SetLogFormatter 日志格式化器
	SetLogFormatter(func(entry *LogEntity) string)
	// ShortLogEnable 是否使用短日志（true则只包含调用者的相对路径）
	ShortLogEnable(shortLog bool)
	// ConsoleEnable 是否允许控制台输出
	ConsoleEnable(console bool)
	// ColorEnable 是否需要彩色输出
	ColorEnable(color bool)
	// Destroy 销毁
	Destroy()
}

// LogEntity
// @Description: 日志消息体
// @Data 2023-02-27 10:07:03
type LogEntity struct {
	LogTime  time.Time //日志时间
	LogLevel LogLevel  //日志级别
	LogFile  string    //产生日志的文件
	LineNum  int       //行号
	Msg      string    // 日志内容
}
