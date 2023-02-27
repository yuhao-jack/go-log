## go-log

> 是一个简单明了的日志框架。我们可以使用它来替换不易使用的内置日志框架。这一框架正在不断改进中，如有疑问可Email:154826195@qq.com

#### 快速开始

#### 安装

```
go get github.com/yuhao-jack/go-log
```

#### example

```
func TestDemo1(t *testing.T) {
	logger := go_log.NewGoLog(&go_log.GoLogConfig{
		LogLevel:       go_log.LoglevelDebug, //指定日志级别，小于该级别日志不处理
		ShortLogEnable: true,//使用短日志（推荐）
		MsgChan:        make(chan string, 256),//缓冲长度
		Writer:         nil,
		ConsoleEnable:  true,//控制台输出
		ColorEnable:    true,//带颜色输出
	})
	defer logger.Destroy()//阻塞直到MsgChan中的消息消费完
	logger.Debug("hello world")
	logger.ColorEnable(false)
	logger.Info("hello world")
	logger.SetLogFormatter(logFormatter)
	logger.Warn("hello world")
	file, err := os.Create("test.log") //创建文件
	if err != nil {
		fmt.Println("create file failed,err:", err)
		os.Exit(1)
	}
	logger.SetLohWriter(file)//日志落盘
	logger.Error("hello world")
}

func logFormatter(entry *go_log.LogEntity) string {
	bytes, err := json.Marshal(entry)
	if err != nil {
		return "\n"
	}
	return string(bytes) + "\n"
}

```

