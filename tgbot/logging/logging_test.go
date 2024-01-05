package logging_test

import (
	"strings"
	"testing"

	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestFiltering(t *testing.T) {

	core, observed := observer.New(zap.DebugLevel)
	baseLogger := logging.FromCore(core)
	loggers := logging.NewLoggersDefault(baseLogger)

	loggers.SetFilter(func(e zapcore.Entry, f []zapcore.Field) bool {
		return strings.Contains(e.LoggerName, "er.bl")
	})

	loggers.IncreaseLevel(zap.InfoLevel)

	loggers.ChatsDispatcher().Named("blah").Info("info")
	loggers.ChatsDispatcher().Named("blah").Debug("debug")
	loggers.Tgbot().Debug("debug")
	loggers.Tgbot().Info("info")

	assert.Equal(t, []observer.LoggedEntry{
		{Entry: zapcore.Entry{Message: "info", LoggerName: "ChatsDispatcher.blah"}, Context: []zapcore.Field{}},
	}, observed.AllUntimed())

}
