package go_log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GoLogConfig
// @Description:GoLog 配置类
type GoLogConfig struct {
	LogLevel       LogLevel    `json:"log_level"`        //日志级别
	ShortLogEnable bool        `json:"short_log_enable"` //是否使用短日志
	MsgChan        chan string `json:"msg_chan"`         //消息管道（缓冲区）
	Writer         io.Writer   `json:"-"`                //输出流 可以使用文件、网络
	ConsoleEnable  bool        `json:"console_enable"`   //控制台输出
	ColorEnable    bool        //颜色输出
}

// GoLog
// @Description: GoLog 实体类
type GoLog struct {
	sync.RWMutex
	logLevel       LogLevel                      //日志级别
	shortLogEnable bool                          //是否使用短日志
	msgChan        chan string                   //消息管道（缓冲区）
	writer         io.Writer                     //输出流
	consoleEnable  bool                          //控制台输出
	colorEnable    bool                          //颜色输出
	waiter         sync.WaitGroup                //阻塞
	logFormatter   func(entry *LogEntity) string //格式化器
}

// DefaultGoLog
//
//	@Description: 根据默认配置创建一个对象实例
//	@Author yuhao
//	@Data 2023-02-27 14:25:54
//	@return *GoLog
func DefaultGoLog() *GoLog {
	g := &GoLog{
		RWMutex:        sync.RWMutex{},
		logLevel:       LoglevelInfo,
		shortLogEnable: true,
		msgChan:        make(chan string, 256),
		writer:         nil,
		consoleEnable:  true,
		colorEnable:    true,
		waiter:         sync.WaitGroup{},
	}
	go g.consumeMsgChan()
	return g
}

var once = sync.Once{}
var singleGoLog *GoLog

// GetSingleGoLog
//
//	@Description: 获取单例GoLog实例
//	@Author yuhao
//	@Data 2023-02-27 14:29:10
//	@return *GoLog
func GetSingleGoLog() *GoLog {
	once.Do(func() {
		singleGoLog = DefaultGoLog()
	})
	return singleGoLog
}

// NewGoLog
//
//	@Description: 创建日志
//	@param config
//	@return *GoLog
func NewGoLog(config *GoLogConfig) ILogger {
	g := &GoLog{
		RWMutex:        sync.RWMutex{},
		logLevel:       config.LogLevel,
		shortLogEnable: config.ShortLogEnable,
		msgChan:        config.MsgChan,
		writer:         nil,
		consoleEnable:  config.ConsoleEnable,
		colorEnable:    config.ColorEnable,
		waiter:         sync.WaitGroup{},
	}

	go g.consumeMsgChan()
	return g
}
func (g *GoLog) Trace(format string, msg ...any) {
	if g.logLevel.LevelNum() > LoglevelTrace.LevelNum() {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprintf(format, msg...)
		entity := LogEntity{
			LogTime:  time.Now(),
			LogLevel: LoglevelTrace,
			LogFile:  g.fileIdx(file),
			LineNum:  line,
			Msg:      data,
		}
		if g.logFormatter != nil {
			g.msgChan <- g.logFormatter(&entity)
		} else {
			g.msgChan <- g.formatMsg(&entity)
		}
	}
}

func (g *GoLog) Debug(format string, msg ...any) {
	if g.logLevel.LevelNum() > LoglevelDebug.LevelNum() {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprintf(format, msg...)
		entity := LogEntity{
			LogTime:  time.Now(),
			LogLevel: LoglevelDebug,
			LogFile:  g.fileIdx(file),
			LineNum:  line,
			Msg:      data,
		}
		if g.logFormatter != nil {
			g.msgChan <- g.logFormatter(&entity)
		} else {
			g.msgChan <- g.formatMsg(&entity)
		}
	}
}

func (g *GoLog) Info(format string, msg ...any) {
	if g.logLevel.LevelNum() > LoglevelInfo.LevelNum() {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprintf(format, msg...)
		entity := LogEntity{
			LogTime:  time.Now(),
			LogLevel: LoglevelInfo,
			LogFile:  g.fileIdx(file),
			LineNum:  line,
			Msg:      data,
		}
		if g.logFormatter != nil {
			g.msgChan <- g.logFormatter(&entity)
		} else {
			g.msgChan <- g.formatMsg(&entity)
		}
	}
}

func (g *GoLog) Warn(format string, msg ...interface{}) {
	if g.logLevel.LevelNum() > LoglevelWarn.LevelNum() {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprintf(format, msg...)
		entity := LogEntity{
			LogTime:  time.Now(),
			LogLevel: LoglevelWarn,
			LogFile:  g.fileIdx(file),
			LineNum:  line,
			Msg:      data,
		}
		if g.logFormatter != nil {
			g.msgChan <- g.logFormatter(&entity)
		} else {
			g.msgChan <- g.formatMsg(&entity)
		}
	}
}

func (g *GoLog) Error(format string, msg ...interface{}) {
	if g.logLevel.LevelNum() > LoglevelError.LevelNum() {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprintf(format, msg...)
		entity := LogEntity{
			LogTime:  time.Now(),
			LogLevel: LoglevelError,
			LogFile:  g.fileIdx(file),
			LineNum:  line,
			Msg:      data,
		}
		if g.logFormatter != nil {
			g.msgChan <- g.logFormatter(&entity)
		} else {
			g.msgChan <- g.formatMsg(&entity)
		}
	}
}

func (g *GoLog) SetLogLevel(loglevel LogLevel) {
	g.RLock()
	defer g.RUnlock()
	g.logLevel = loglevel
}

func (g *GoLog) SetLohWriter(writer io.Writer) {
	g.RLock()
	defer g.RUnlock()
	g.writer = writer
}

func (g *GoLog) SetLogFormatter(f func(entry *LogEntity) string) {
	g.RLock()
	defer g.RUnlock()
	g.logFormatter = f
}

func (g *GoLog) ShortLogEnable(shortLog bool) {
	g.RLock()
	defer g.RUnlock()
	g.shortLogEnable = shortLog
}

func (g *GoLog) ConsoleEnable(console bool) {
	g.RLock()
	defer g.RUnlock()
	g.consoleEnable = console
}

func (g *GoLog) ColorEnable(color bool) {
	g.RLock()
	defer g.RUnlock()
	g.colorEnable = color
}

func (g *GoLog) Destroy() {
	close(g.msgChan)
	g.waiter.Wait()

}

// formatMsg
//
//	@Description: 格式化日志明细
//	@receiver g
//	@param level
//	@param file
//	@param line
//	@param msg
//	@return string
func (g *GoLog) formatMsg(entry *LogEntity) string {
	var detail string
	if g.colorEnable {
		detail = fmt.Sprint(
			Cyan.WithColorEnd(entry.LogTime.Format(string(DefaultLayout))),
			fmt.Sprintf("%18s", " ["+Green.WithColorEnd(string(entry.LogLevel))+"] "),
			fmt.Sprintf("%30s", entry.LogFile+":"+strconv.Itoa(entry.LineNum)+" \t:"),
			entry.Msg,
		)
	} else {
		detail = fmt.Sprint(
			entry.LogTime.Format(string(DefaultLayout)),
			fmt.Sprintf("%18s", " ["+entry.LogLevel+"] "),
			fmt.Sprintf("%30s", entry.LogFile+":"+strconv.Itoa(entry.LineNum)+" \t:"),

			entry.Msg,
		)
	}

	return detail + "\n"
}

// fileIdx
//
//	@Description: 获取文件地址 如果是长文件则直接返回
//	@receiver g
//	@param file 文件绝对地址
//	@return string 文件地址
func (g *GoLog) fileIdx(file string) string {
	if !g.shortLogEnable {
		return file
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file
}

// consumeMsgChan
//
//	@Description: 消费消息管道的消息
//	@receiver g
func (g *GoLog) consumeMsgChan() {
	g.waiter.Add(1)
	for {
		select {
		case msg, ok := <-g.msgChan:
			if !ok { //此时说明管道已经关闭
				g.waiter.Done()
				return
			}
			if g.consoleEnable {
				_, _ = os.Stdout.WriteString(msg)
			}
			if g.writer != nil {
				_, _ = g.writer.Write([]byte(msg))
			}
		}
	}
}
