package logging

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogsSystemProps struct {
	TeeUserToConsole bool

	EnableFileLog     bool
	EnableUserFileLog bool
}

type LogsSystem struct {
	baseLogger *zap.Logger
	logsFolder string

	props LogsSystemProps

	logSystemLogger *zap.Logger

	logsPrefix string
}

func NewLogsSystem(base *zap.Logger, logsFolder string, props LogsSystemProps) (*LogsSystem, error) {

	logsPrefix := time.Now().Format("2006-01-02_15-04-05")

	ls := &LogsSystem{
		baseLogger:      base,
		logSystemLogger: base.Named(LoggerNameLogsSystem),
		logsFolder:      logsFolder,
		props:           props,
		logsPrefix:      logsPrefix,
	}

	if props.EnableFileLog {
		err := ls.attachMainFileLogger()

		if err != nil {
			return nil, err
		}
	}

	return ls, nil
}

func (ul *LogsSystem) attachMainFileLogger() error {

	mainLogFile, err := ul.openMainLogFile()

	if err != nil {
		ul.logSystemLogger.Error("failed to open log file", zap.Error(err))
		return err
	}

	ul.baseLogger = ul.baseLogger.WithOptions(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {

			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(zapcore.Lock(mainLogFile)),
				zap.DebugLevel,
			)

			if !ul.props.TeeUserToConsole {
				return fileCore
			}

			return zapcore.NewTee(core, fileCore)
		}),
	)

	return nil
}

func (ul *LogsSystem) openMainLogFile() (*os.File, error) {
	fname := fmt.Sprintf("%s_main.log", ul.logsPrefix)
	mainLogFile := path.Join(ul.logsFolder, fname)
	ul.logSystemLogger.Info("Opening log file", zap.String("path", mainLogFile))

	file, err := os.OpenFile(mainLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		ul.logSystemLogger.Error("failed to open log file", zap.Error(err))
		return nil, err
	}

	return file, nil

}

func (ul *LogsSystem) openUserLogFile(userID int64) (*os.File, error) {
	fname := fmt.Sprintf("%s_user_%d.log", ul.logsPrefix, userID)

	userLogFile := path.Join(ul.logsFolder, fname)
	ul.logSystemLogger.Info("Opening log file", zap.String("path", userLogFile))

	return os.OpenFile(userLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (ul *LogsSystem) ApplicationChatLogger(
	l *zap.Logger,
	tc *telegram.TelegramUpdateContext,
) *zap.Logger {

	l = l.Named(LoggerNameApplicationChat)

	if !ul.props.EnableUserFileLog {
		return l
	}

	file, err := ul.openUserLogFile(tc.ChatID)

	if err != nil {
		ul.logSystemLogger.Error("failed to open log file", zap.Error(err))
		return l
	}

	l = l.WithOptions(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {

			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(zapcore.Lock(file)),
				zap.DebugLevel,
			)

			if !ul.props.TeeUserToConsole {
				return fileCore
			}

			return zapcore.NewTee(core, fileCore)
		}),
	)

	return l
}

func (ul *LogsSystem) Loggers() Loggers {
	return NewLoggersDefault(ul.baseLogger)
}
