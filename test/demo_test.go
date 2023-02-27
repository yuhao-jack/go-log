package test

import (
	"encoding/json"
	"fmt"
	go_log "github.com/yuhao-jack/go-log"

	"os"
	"testing"
	"time"
)

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

func TestDemo2(t *testing.T) {
	now := time.Now().Unix()
	fmt.Println(time.Unix(now, 0).Format(string(go_log.DefaultLayout)))
	fiveMinute := int64(5 * 60)
	fmt.Println(time.Unix(now/fiveMinute*fiveMinute, 0).Format(string(go_log.DefaultLayout)))

}
