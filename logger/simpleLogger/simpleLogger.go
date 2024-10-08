package simplelogger

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gnanasuryateja/golib/constants"
	"github.com/gnanasuryateja/golib/logger"
)

type SimpleLoggerParams struct {
	ServiceName          string  // ServiceName is the name of the service in which you are working |
	LogLevel             *string // LogLevel is the log level configured from env |
	SkipLevelForFuncInfo *int    // SkipLevelForFuncInfo refers to the skip param to pass in runtime.Caller(skip)
	Env                  string  // Env is the environment in which the application is running
}

type simpleLogger struct {
	ServiceName          string // ServiceName is the name of the service in which you are working |
	LogLevel             string // LogLevel is the log level configured from env |
	SkipLevelForFuncInfo int    // SkipLevelForFuncInfo refers to the skip param to pass in runtime.Caller(skip)
	Env                  string // Env is the environment in which the application is running
}

func (sl simpleLogger) validate() error {
	if sl.ServiceName == "" {
		return fmt.Errorf("service name is passed as empty")
	}
	if sl.LogLevel != "" {
		if !(strings.EqualFold(sl.LogLevel, constants.LOG_LEVEL_DEBUG) ||
			strings.EqualFold(sl.LogLevel, constants.LOG_LEVEL_ERROR) ||
			strings.EqualFold(sl.LogLevel, constants.LOG_LEVEL_INFO) ||
			strings.EqualFold(sl.LogLevel, constants.LOG_LEVEL_WARN)) {
			return fmt.Errorf("invalid log level... %s is not supported by simpleLogger", sl.LogLevel)
		}
	}
	return nil
}

func NewSimpleLogger(loggerParams SimpleLoggerParams) (logger.Logger, error) {
	var simpleLogger simpleLogger
	simpleLogger.ServiceName = loggerParams.ServiceName
	if loggerParams.LogLevel == nil || *loggerParams.LogLevel == "" {
		simpleLogger.LogLevel = constants.LOG_LEVEL_DEBUG
	} else {
		simpleLogger.LogLevel = *loggerParams.LogLevel
	}
	if loggerParams.SkipLevelForFuncInfo == nil {
		simpleLogger.SkipLevelForFuncInfo = 2
	} else {
		simpleLogger.SkipLevelForFuncInfo = *loggerParams.SkipLevelForFuncInfo
	}
	simpleLogger.Env = loggerParams.Env
	err := simpleLogger.validate()
	if err != nil {
		return nil, err
	}
	return simpleLogger, nil
}

// Debug logs would be printed only when LogLevel is set to DEBUG |
func (sl simpleLogger) Debug(ctx context.Context, message string) {
	funcName, fileName, lineNo := logger.GetCurrentFuncInfo(sl.SkipLevelForFuncInfo)
	if strings.EqualFold(sl.LogLevel, "DEBUG") {
		fServiceName := fmt.Sprintf("[%v]", sl.ServiceName)
		fmt.Println(buildSimpleLog(fServiceName, "debug", funcName, fileName, lineNo, message))
	}
}

// Error logs would always be printed |
func (sl simpleLogger) Error(ctx context.Context, err error) {
	funcName, fileName, lineNo := logger.GetCurrentFuncInfo(sl.SkipLevelForFuncInfo)
	fServiceName := fmt.Sprintf("[%v]", sl.ServiceName)
	fmt.Println(buildSimpleLog(fServiceName, "error", funcName, fileName, lineNo, err.Error()))
}

// Info logs would always be printed |
func (sl simpleLogger) Info(ctx context.Context, message string) {
	funcName, fileName, lineNo := logger.GetCurrentFuncInfo(sl.SkipLevelForFuncInfo)
	fServiceName := fmt.Sprintf("[%v]", sl.ServiceName)
	fmt.Println(buildSimpleLog(fServiceName, "info", funcName, fileName, lineNo, message))
}

// Warn logs would be printed only when LogLevel is set to DEBUG |
func (sl simpleLogger) Warn(ctx context.Context, message string) {
	funcName, fileName, lineNo := logger.GetCurrentFuncInfo(sl.SkipLevelForFuncInfo)
	if strings.EqualFold(sl.LogLevel, "DEBUG") {
		fServiceName := fmt.Sprintf("[%v]", sl.ServiceName)
		fmt.Println(buildSimpleLog(fServiceName, "warn", funcName, fileName, lineNo, message))
	}
}

func buildSimpleLog(serviceName string, logFunc string, funcName string, fileName string, lineNo int, logMsg string) string {
	return serviceName + " [" + time.Now().UTC().Format(constants.SIMPLE_LOGGER_TIME_FORMAT) + "] " + strings.ToUpper(logFunc) + ": " + funcName + "() " + fileName + fmt.Sprintf(":%v ", lineNo) + logMsg
}
