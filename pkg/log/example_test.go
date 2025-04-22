package log_test

import (
	"os"

	"github.com/soulchildtc/mcp-server-weather/pkg/log"
)

func Example() {
	// 使用默认日志记录器
	log.Info("这是一条信息日志")
	log.Warn("这是一条警告日志")
	log.Error("这是一条错误日志")

	// 设置日志级别
	log.SetLevel(log.DebugLevel)
	log.Debug("设置日志级别后可以看到调试日志")

	// 创建自定义日志记录器
	logger := log.NewLogger(log.Options{
		Level:      log.InfoLevel,
		Output:     os.Stdout,
		UseColor:   true,
		TimeFormat: "2006/01/02 15:04:05",
		CallerSkip: 1,
	})

	// 使用自定义日志记录器
	logger.Info("使用自定义日志记录器")
	logger.SetLevel(log.WarnLevel)
	logger.Info("这条信息不会显示，因为级别低于警告级别")
	logger.Warn("这条警告会显示")

	// 关闭颜色
	logger.ToggleColor(false)
	logger.Error("这条错误日志不会使用颜色")

	// 创建一个写入文件的日志记录器
	// (实际应用中应该检查错误)
	file, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	fileLogger := log.NewLogger(log.Options{
		Level:    log.InfoLevel,
		Output:   file,
		UseColor: false, // 文件日志通常不使用颜色
	})

	fileLogger.Info("这条日志会写入到文件")
	fileLogger.Error("错误信息也会写入到文件")

	// Output:
	// 示例输出 (实际输出会包含时间戳和文件位置)
}

// 实际应用示例
func ExampleUsage() {
	// 在应用程序初始化时配置日志
	log.SetLevel(log.DebugLevel)

	// 在HTTP服务器中使用
	// router.GET("/api/users", func(c *gin.Context) {
	// 	log.Info("收到用户列表请求")
	// 	// ...处理请求
	// 	if err != nil {
	// 		log.Error("获取用户列表失败: %v", err)
	// 		// ...返回错误
	// 		return
	// 	}
	// 	log.Debug("返回用户列表: %v", users)
	// 	// ...返回结果
	// })
}
