package go_log

import (
	"io"
	"time"
)

type LogLevel uint8 //日志级别

type TimeLayOut string //时间格式

type Color string //颜色

const (
	LoglevelDebug LogLevel = 1
	LoglevelInfo  LogLevel = 2
	LoglevelWarn  LogLevel = 3
	LoglevelError LogLevel = 4
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

// DefaultLayout 默认日志时间格式
// 2006为Golang诞生时间，15是下午3点。帮助记忆的方法：1月2日3点4分5秒，2006年，-7时区，正好是1234567
const DefaultLayout TimeLayOut = "2006-01-02 15:04:05.000"

// String
//
//	@Description: 获取日志级别字符串
//	@receiver l
//	@return string
func (l LogLevel) String() string {
	switch l {
	case LoglevelDebug:
		return "DEBUG"
	case LoglevelInfo:
		return "INFO"
	case LoglevelWarn:
		return "WARN"
	case LoglevelError:
		return "ERROR"
	default:
		return ""
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
	// Debug Debug级别日志
	Debug(msg ...any)
	// Info Info级别日志
	Info(msg ...any)
	// Warn Warn级别日志
	Warn(msg ...any)
	// Error Error级别日志
	Error(msg ...any)
	// SetLogLevel 设置日志级别
	SetLogLevel(loglevel LogLevel)
	// SetLohWriter 设置输出流
	SetLohWriter(writer io.Writer)
	// ShortLogEnable 是否使用短日志（true则只包含调用者的相对路径）
	ShortLogEnable(shortLog bool)
	// ConsoleEnable 是否允许控制台输出
	ConsoleEnable(console bool)
	// ColorEnable 是否需要彩色输出
	ColorEnable(color bool)
	// Destroy 销毁
	Destroy()
}
