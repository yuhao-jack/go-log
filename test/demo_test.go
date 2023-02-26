package test

import (
	go_log "go-log"
	"testing"
)

func TestDemo(t *testing.T) {
	logger := go_log.NewGoLog(&go_log.GoLogConfig{
		LogLevel:       go_log.LoglevelDebug,
		ShortLogEnable: true,
		MsgChan:        make(chan string, 256),
		Writer:         nil,
		ConsoleEnable:  true,
	})
	defer logger.Destroy()
	logger.Debug("hello world")
}
