package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level 定义日志级别
type Level int

const (
	// 日志级别定义
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var (
	// 级别对应的字符串
	levelNames = map[Level]string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		ErrorLevel: "ERROR",
		FatalLevel: "FATAL",
	}

	// 级别对应的颜色
	levelColors = map[Level]string{
		DebugLevel: "\033[37m", // 白色
		InfoLevel:  "\033[32m", // 绿色
		WarnLevel:  "\033[33m", // 黄色
		ErrorLevel: "\033[31m", // 红色
		FatalLevel: "\033[35m", // 紫色
	}

	// 重置颜色
	resetColor = "\033[0m"

	// 默认日志记录器
	defaultLogger *Logger
	once          sync.Once
)

// Logger 日志记录器结构体
type Logger struct {
	level      Level
	output     io.Writer
	logger     *log.Logger
	mu         sync.Mutex
	useColor   bool
	timeFormat string
	callerSkip int
}

// Options 日志配置选项
type Options struct {
	// 日志级别
	Level Level
	// 日志输出位置
	Output io.Writer
	// 是否使用颜色
	UseColor bool
	// 时间格式
	TimeFormat string
	// 调用者深度
	CallerSkip int
}

// 初始化默认日志记录器
func init() {
	defaultLogger = NewLogger(Options{
		Level:      InfoLevel,
		Output:     os.Stdout,
		UseColor:   true,
		TimeFormat: "2006-01-02 15:04:05.000",
		CallerSkip: 4,
	})
}

// NewLogger 创建新的日志记录器
func NewLogger(opts Options) *Logger {
	// 设置默认值
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = "2006-01-02 15:04:05.000"
	}
	if opts.CallerSkip == 0 {
		opts.CallerSkip = 3
	}

	return &Logger{
		level:      opts.Level,
		output:     opts.Output,
		logger:     log.New(opts.Output, "", 0),
		useColor:   opts.UseColor,
		timeFormat: opts.TimeFormat,
		callerSkip: opts.CallerSkip,
	}
}

// formatHeader 格式化日志头部
func (l *Logger) formatHeader(level Level, calldepth int) string {
	now := time.Now().Format(l.timeFormat)
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	// 提取文件名
	file = filepath.Base(file)

	// 格式化日志头
	levelName := levelNames[level]
	if l.useColor {
		levelColor := levelColors[level]
		return fmt.Sprintf("%s [%s%5s%s] %s:%d", now, levelColor, levelName, resetColor, file, line)
	}
	return fmt.Sprintf("%s [%5s] %s:%d", now, levelName, file, line)
}

// log 记录日志的内部方法
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	header := l.formatHeader(level, l.callerSkip)
	var msg string
	if format == "" {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	// 输出日志
	l.logger.Println(header, msg)

	// 如果是致命错误，则退出程序
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug 输出调试级别日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DebugLevel, format, args...)
}

// Info 输出信息级别日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(InfoLevel, format, args...)
}

// Warn 输出警告级别日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WarnLevel, format, args...)
}

// Error 输出错误级别日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ErrorLevel, format, args...)
}

// Fatal 输出致命错误日志并退出程序
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FatalLevel, format, args...)
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput 设置日志输出位置
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = w
	l.logger = log.New(w, "", 0)
}

// ToggleColor 开关颜色输出
func (l *Logger) ToggleColor(useColor bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.useColor = useColor
}

// 以下是包级别的函数，使用默认logger

// Debug 输出调试级别日志
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info 输出信息级别日志
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn 输出警告级别日志
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error 输出错误级别日志
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// Fatal 输出致命错误日志并退出程序
func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

// SetLevel 设置默认日志记录器的日志级别
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// SetOutput 设置默认日志记录器的输出
func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

// ToggleColor 开关默认日志记录器的颜色输出
func ToggleColor(useColor bool) {
	defaultLogger.ToggleColor(useColor)
}

// 获取单例默认记录器
func GetDefaultLogger() *Logger {
	return defaultLogger
}

// 简单日志方法
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}
