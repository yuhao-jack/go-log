package go_log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// GoLogConfig
// @Description:GoLog 配置类，当RollLogByTime、RollLogBySize二者都不为空时只会生效一个，优选使用RollLogByTime
type GoLogConfig struct {
	LogLevel       LogLevel    `json:"log_level"`        //日志级别
	ShortLogEnable bool        `json:"short_log_enable"` //是否使用短日志
	MsgChan        chan string `json:"msg_chan"`         //消息管道（缓冲区）
	Writer         io.Writer   `json:"-"`                //输出流 可以使用文件、网络
	ConsoleEnable  bool        `json:"console_enable"`   //控制台输出
	ColorEnable    bool        `json:"color_enable"`     //颜色输出
	LogDir         string      `json:"log_dir"`          //日志存放目录
	LogName        string      `json:"log_name"`         //日志文件名
	RollLogByTime  string      `json:"roll_log_by_time"` //根据时间滚动 如:5m表示五分钟滚动一个，为了便于管理这里会把时间整块分，如16:56:23则会写进16:55:00这个时间块的文件中
	RollLogBySize  int64       `json:"roll_log_by_size"` //根据文件大小滚动，单位KB，

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
	logDir         string                        //日志存放目录
	logName        string                        //日志文件名
	rollLogByTime  time.Duration                 //根据时间滚动 如:5m表示五分钟滚动一个，为了便于管理这里会把时间整块分，如16:56:23则会写进16:55:00这个时间块的文件中
	rollLogBySize  int64                         //根据文件大小滚动，单位KB，
	logFile        *os.File                      //日志文件句柄
	lastTimeBlock  string                        //文件最后变更时间的时间块
	logFileSize    int64                         //当前日志文件的大小
	compressChan   chan string                   //压缩文件信号管道，将要压缩的文件名丢入管道
	closeFlag      bool
}

// DefaultGoLog
//
//	@Description: 根据默认配置创建一个对象实例 Info 级别
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
//	@Description: 获取单例GoLog实例 Info 级别
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
		logDir:         config.LogDir,
		logName:        config.LogName,
		rollLogBySize:  config.RollLogBySize,
	}
	if config.RollLogByTime != "" {
		duration, err := time.ParseDuration(config.RollLogByTime)
		if err != nil {
			panic("Invalid time:" + config.RollLogByTime)
		}
		g.rollLogByTime = duration
	}
	go g.consumeMsgChan()
	return g
}
func (g *GoLog) Trace(format string, msg ...any) {
	if g.logLevel.LevelNum() > LoglevelTrace.LevelNum() || g.closeFlag {
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
	if g.logLevel.LevelNum() > LoglevelDebug.LevelNum() || g.closeFlag {
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
	if g.logLevel.LevelNum() > LoglevelInfo.LevelNum() || g.closeFlag {
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
	if g.logLevel.LevelNum() > LoglevelWarn.LevelNum() || g.closeFlag {
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
	if g.logLevel.LevelNum() > LoglevelError.LevelNum() || g.closeFlag {
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

func (g *GoLog) SetLogDir(logDir string) {
	g.RLock()
	defer g.RUnlock()
	g.logDir = logDir
}

func (g *GoLog) setLogName(logName string) {
	g.RLock()
	defer g.RUnlock()
	g.logName = logName
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
	g.closeFlag = true
	close(g.msgChan)
	close(g.compressChan)
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
		var color Color
		switch entry.LogLevel {
		case LoglevelTrace:
			color = Blue
		case LoglevelDebug:
			color = Magenta
		case LoglevelInfo:
			color = Green
		case LoglevelWarn:
			color = Yellow
		case LoglevelError:
			color = Red
		default:
			color = White
		}
		detail = fmt.Sprint(
			Cyan.WithColorEnd(entry.LogTime.Format(string(DefaultLayout))),
			fmt.Sprintf("%18s", " ["+color.WithColorEnd(string(entry.LogLevel))+"] "),
			fmt.Sprintf("%30s", entry.LogFile+":"+strconv.Itoa(entry.LineNum)+":\t"),
			entry.Msg,
		)
	} else {
		detail = fmt.Sprint(
			entry.LogTime.Format(string(DefaultLayout)),
			fmt.Sprintf("%18s", " ["+entry.LogLevel+"] "),
			fmt.Sprintf("%30s", entry.LogFile+":"+strconv.Itoa(entry.LineNum)+":\t"),
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
	//  目录、文件名不为空 切没有结尾斜杠
	if g.logDir != "" && g.logName != "" && !(strings.HasSuffix(g.logDir, "/") || strings.HasSuffix(g.logDir, "\\")) {
		g.SetLogDir(g.logDir + "/")
		g.setLogName(g.logDir + g.logName)
	}
	if g.compressChan == nil {
		g.compressChan = make(chan string, 2)
	}
	if g.rollLogByTime != 0 || g.rollLogBySize != 0 {
		go g.compressLogFile()

	}
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
			if g.logDir == "" || g.logName == "" {
				continue
			}
			file := g.getLogFile()
			if file == nil {
				continue
			}
			n, err := file.WriteString(msg)
			if err != nil {
				_, _ = os.Stderr.WriteString("write log to " + g.logName + " failed,err:" + err.Error() + "\tdata:" + msg)
			}
			g.logFileSize += int64(n)
		}
	}
}

// compressLogFile
//
//	@Description: 异步压缩文件
//	@receiver g
//	@Author yuhao
//	@Data 2023-02-28 11:30:47
func (g *GoLog) compressLogFile() {
	g.waiter.Add(1)
	for {
		select {
		case s, ok := <-g.compressChan:
			if !ok {
				g.waiter.Done()
				return
			}
			file, err := os.Open(s)
			if err != nil {
				_, _ = os.Stderr.WriteString("open file " + s + " failed,err:" + err.Error())
				continue
			}
			err = Compress([]*os.File{file}, s+".zip")
			if err != nil {
				_, _ = os.Stderr.WriteString("Compress file " + s + ".zip" + " failed,err:" + err.Error())
			}
			_ = os.Remove(s)
		}
	}

}

// getLogFile
//
//	@Description: 获取文件句柄
//	@receiver g
//	@return *os.File
func (g *GoLog) getLogFile() *os.File {
	fileInfo, err := os.Stat(g.logName)
	if os.IsNotExist(err) { //文件不存在
		file, err := os.Create(g.logName)
		if err != nil {
			_, _ = os.Stderr.WriteString("create logfile " + g.logName + " failed,err:" + err.Error())
			return nil
		}
		g.logFile = file
		return file
	}
	//  在同一个时间块但是还没打开
	if g.logFile == nil {
		file, err := os.OpenFile(g.logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			_, _ = os.Stderr.WriteString("open logfile " + g.logName + " failed,err:" + err.Error())
			return nil
		}
		g.logFile = file
	}
	//  文件存在 根据时间滚动文件
	if g.rollLogByTime != 0 {
		return g.getFileByTime(fileInfo)
	}
	//  文件存在 根据文件的大小滚动文件
	if g.rollLogBySize != 0 {
		return g.getFileBySize(fileInfo)
	}

	return g.logFile
}

// getFileByTime
//
//	@Description: 根据时间滚动文件
//	@receiver g
//	@param fileInfo
//	@return *os.File
func (g *GoLog) getFileByTime(fileInfo os.FileInfo) *os.File {
	now := time.Now().Unix()
	duration := int64(g.rollLogByTime.Seconds())
	format := time.Unix(now/duration*duration, 0).Format(string(DateTimeLayout4))
	if g.lastTimeBlock == "" {
		g.lastTimeBlock = fileInfo.ModTime().Format(string(DateTimeLayout4))
	}
	//  不在同一时间块
	if g.lastTimeBlock != format {
		// 如果文件被打开需要关闭
		if g.logFile != nil {
			_ = g.logFile.Close()
			g.logFile = nil
		}
		err := os.Rename(g.logName, g.logName+"-"+g.lastTimeBlock)
		if err != nil {
			_, _ = os.Stderr.WriteString("Rename " + g.logName + " failed,err:" + err.Error())
			return nil
		}
		g.compressChan <- g.logName + "-" + g.lastTimeBlock
		file, err := os.Create(g.logName)
		if err != nil {
			_, _ = os.Stderr.WriteString("create logfile " + g.logName + " failed,err:" + err.Error())
			return nil
		}
		g.lastTimeBlock = format
		g.logFile = file
		return file
	}

	return g.logFile
}

// getFileBySize
//
//	@Description: 根据文件大小滚动文件
//	@receiver g
//	@param fileInfo
//	@return *os.File
func (g *GoLog) getFileBySize(fileInfo os.FileInfo) *os.File {
	sizeKB := fileInfo.Size() / 1024
	// 文件大小超过滚动的大小了需要重命名滚动
	if g.rollLogBySize < sizeKB {
		// 如果文件被打开需要关闭
		if g.logFile != nil {
			_ = g.logFile.Close()
			g.logFile = nil
		}
		fileInfos, err := ioutil.ReadDir(g.logDir)
		if err != nil {
			_, _ = os.Stderr.WriteString("ReadDir " + g.logDir + " failed,err:" + err.Error())
			return nil
		}
		cnt := 1 //获取文件夹中已经存在多少logName文件了
		for _, info := range fileInfos {
			if strings.HasPrefix(info.Name(), g.logName) && strings.HasSuffix(info.Name(), ".zip") {
				cnt++
			}
		}

		err = os.Rename(g.logName, g.logName+"-"+strconv.Itoa(cnt))
		if err != nil {
			_, _ = os.Stderr.WriteString("Rename " + g.logName + " failed,err:" + err.Error())
			return nil
		}

		g.compressChan <- g.logName + "-" + strconv.Itoa(cnt)
		file, err := os.Create(g.logName)
		if err != nil {
			_, _ = os.Stderr.WriteString("create logfile " + g.logName + " failed,err:" + err.Error())
			return nil
		}
		g.logFile = file
		return file
	}
	return g.logFile
}
