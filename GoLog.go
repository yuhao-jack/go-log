package go_log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// GoLogConfig
// @Description:GoLog 配置类
type GoLogConfig struct {
	LogLevel       LogLevel    `json:"log_level"`        //日志级别
	ShortLogEnable bool        `json:"short_log_enable"` //是否使用短日志
	MsgChan        chan string `json:"msg_chan"`         //消息管道（缓冲区）
	Writer         io.Writer   `json:"-"`                //输出流 可以使用文件、网络
	ConsoleEnable  bool        `json:"console_enable"`   //控制台输出
}

// GoLog
// @Description: GoLog 实体类
type GoLog struct {
	sync.RWMutex
	logLevel       LogLevel       //日志级别
	shortLogEnable bool           //是否使用短日志
	msgChan        chan string    //消息管道（缓冲区）
	writer         io.Writer      //输出流
	consoleEnable  bool           //控制台输出
	waiter         sync.WaitGroup //阻塞
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
		waiter:         sync.WaitGroup{},
	}
	g.waiter.Add(1)
	go g.consumeMsgChan()
	return g
}

func (g *GoLog) Debug(msg ...any) {
	if g.logLevel > LoglevelDebug {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprint(msg)
		g.msgChan <- g.formatMsg(LoglevelDebug.String(), file, strconv.Itoa(line), data[1:len(data)-1])
	}
}

func (g *GoLog) Info(msg ...any) {
	if g.logLevel > LoglevelInfo {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprint(msg)
		g.msgChan <- g.formatMsg(LoglevelInfo.String(), file, strconv.Itoa(line), data[1:len(data)-1])
	}
}

func (g *GoLog) Warn(msg ...interface{}) {
	if g.logLevel > LoglevelWarn {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprint(msg)
		g.msgChan <- g.formatMsg(LoglevelWarn.String(), file, strconv.Itoa(line), data[1:len(data)-1])
	}
}

func (g *GoLog) Error(msg ...interface{}) {
	if g.logLevel > LoglevelError {
		return
	}
	if _, file, line, ok := runtime.Caller(1); ok {
		data := fmt.Sprint(msg)
		g.msgChan <- g.formatMsg(LoglevelError.String(), file, strconv.Itoa(line), data[1:len(data)-1])
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
	g.consoleEnable = color
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
func (g *GoLog) formatMsg(level, file, line, msg string) string {
	detail := fmt.Sprint(Cyan.WithColorEnd(DefaultLayout.String()), " ", level, " ", file, ":", line, " ", msg)
	return detail
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

	for {
		select {
		case msg, ok := <-g.msgChan:
			if !ok { //此时说明管道已经关闭
				g.waiter.Done()
				break
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
