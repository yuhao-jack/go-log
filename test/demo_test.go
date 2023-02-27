package test

import (
	"encoding/json"
	"fmt"
	go_log "github.com/yuhao-jack/go-log"
	"os"
	"testing"
	"time"
)

// TestDemo1
//
//	@Description: 自定义配置
//	@Author yuhao
//	@Data 2023-02-27 16:43:11
//	@param t
func TestDemo1(t *testing.T) {
	logger := go_log.NewGoLog(&go_log.GoLogConfig{
		LogLevel:       go_log.LoglevelDebug,
		ShortLogEnable: true,
		MsgChan:        make(chan string, 256),
		Writer:         nil,
		ConsoleEnable:  true,
		ColorEnable:    true,
	})
	defer logger.Destroy()
	logger.Debug("hello world")
	logger.ColorEnable(false)
	logger.Info("hello world")
	logger.SetLogFormatter(logFormatter)
	logger.Warn("hello world")
	file, err := os.Create("test.log")
	if err != nil {
		fmt.Println("create file failed,err:", err)
		os.Exit(1)
	}
	logger.SetLohWriter(file)
	logger.Error("hello world")

}

func logFormatter(entry *go_log.LogEntity) string {
	bytes, err := json.Marshal(entry)
	if err != nil {
		return "\n"
	}
	return string(bytes) + "\n"
}

// TestDemo2
//
//	@Description: 使用单例的日志对象，任何地方调用返回的示例都相同
//	@param t
func TestDemo2(t *testing.T) {
	singleGoLog := go_log.GetSingleGoLog()
	defer singleGoLog.Destroy()
	singleGoLog.Info("hello world")
}

// TestDemo3
//
//	@Description: 使用默认的日志（DefaultGoLog每次调用将创建一个对象）
//	@Author yuhao
//	@Data 2023-02-27 16:43:29
//	@param t
func TestDemo3(t *testing.T) {
	defaultGoLog := go_log.DefaultGoLog()
	defer defaultGoLog.Destroy()
	defaultGoLog.Info("hello world")
}

// TestDemo4
//
//	@Description: 日志落盘
//	@param t
func TestDemo4(t *testing.T) {
	logger := go_log.NewGoLog(&go_log.GoLogConfig{
		LogLevel:       go_log.LoglevelDebug,
		ShortLogEnable: true,
		MsgChan:        make(chan string, 256),
		Writer:         nil,
		ConsoleEnable:  true,
		ColorEnable:    true,
		LogDir:         "./",
		LogName:        "test1.log",
	})
	defer logger.Destroy()
	for i := 0; i < 10; i++ {
		logger.Info("我的名字叫%s,我今年%d岁了", "二狗子", 18)
	}
}

// TestDemo5
//
//	@Description: 根据时间块滚动文件
//	@param t
func TestDemo5(t *testing.T) {
	logger := go_log.NewGoLog(&go_log.GoLogConfig{
		LogLevel:       go_log.LoglevelDebug,
		ShortLogEnable: true,
		MsgChan:        make(chan string, 256),
		Writer:         nil,
		ConsoleEnable:  true,
		ColorEnable:    true,
		LogDir:         "./",
		LogName:        "test2.log",
		RollLogByTime:  "1m",
	})
	defer logger.Destroy()
	for i := 0; i < 1000; i++ {
		logger.Info("我的名字叫%s,我今年%d岁了", "二狗子", 18)
		time.Sleep(2 * time.Second)

	}
}
func TestDemo6(t *testing.T) {
	logger := go_log.NewGoLog(&go_log.GoLogConfig{
		LogLevel:       go_log.LoglevelDebug,
		ShortLogEnable: true,
		MsgChan:        make(chan string, 256),
		Writer:         nil,
		ConsoleEnable:  true,
		ColorEnable:    true,
		LogDir:         "./",
		LogName:        "test3.log",
		RollLogBySize:  2,
	})
	defer logger.Destroy()
	for i := 0; i < 1000; i++ {
		logger.Info("我的名字叫%s,我今年%d岁了", "二狗子", 18)
		time.Sleep(20 * time.Millisecond)

	}
}
