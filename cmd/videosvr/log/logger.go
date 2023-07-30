package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
)

// 在上下文不重要的环境中，zap中的sugared Logger比 Logger性能更好
var sugarLogger *zap.SugaredLogger
var logger *zap.Logger

func InitLogger() {
	// 获取编码器
	encoderConfig := zap.NewProductionEncoderConfig()            // NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 按级别显示不同颜色，不需要的话取值zapcore.CapitalLevelEncoder就可以了
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder        //显示完整文件路径
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// test.log记录全量日志
	logF := &lumberjack.Logger{
		Filename:   "./log/test.log",
		MaxSize:    200,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   false,
	}
	c1 := zapcore.NewCore(encoder, zapcore.AddSync(logF), zapcore.DebugLevel)

	// test.err.log记录ERROR级别的日志
	errF := &lumberjack.Logger{
		Filename:   "./log/test.err.log",
		MaxSize:    200,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   false,
	}

	c2 := zapcore.NewCore(encoder, zapcore.AddSync(errF), zap.ErrorLevel)
	c3 := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	// 使用NewTee将c1和c2合并到core
	core := zapcore.NewTee(c1, c2, c3)

	// 调用者详细信息, 间接调用log
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger = logger.Sugar()
}

func Test(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

func Infof(s string, v ...interface{}) {
	sugarLogger.Infof(s, v...)
}

func Infow(s string, v ...interface{}) {
	sugarLogger.Infow(s, v...)
}

func Info(v ...interface{}) {
	sugarLogger.Info(v...)
}

func Debugf(s string, v ...interface{}) {
	sugarLogger.Debugf(s, v...)
}

func Debugw(s string, v ...interface{}) {
	sugarLogger.Debugw(s, v...)
}

func Debug(v ...interface{}) {
	sugarLogger.Debug(v...)
}

func Errorf(s string, v ...interface{}) {
	sugarLogger.Errorf(s, v...)
}

func Errorw(s string, v ...interface{}) {
	sugarLogger.Errorw(s, v...)
}

func Error(v ...interface{}) {
	sugarLogger.Error(v...)
}

func Fatalf(s string, v ...interface{}) {
	sugarLogger.Fatalf(s, v...)
}

func Fatalw(s string, v ...interface{}) {
	sugarLogger.Fatalw(s, v...)
}

func Fatal(v ...interface{}) {
	sugarLogger.Error(v...)
}

func Sync() {
	logger.Sync()
}
