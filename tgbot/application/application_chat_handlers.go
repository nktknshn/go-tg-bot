package application

import (
	"github.com/nktknshn/go-tg-bot/tgbot/logging"
	"github.com/nktknshn/go-tg-bot/tgbot/reflection"
	"github.com/nktknshn/go-tg-bot/tgbot/render"
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
	"go.uber.org/zap"
)

func DefaultHandlerCallback[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramContextCallback) {

	logger := ac.Loggers.Handler().With(zap.Int64("UpdateID", tc.UpdateID))

	logger.Info("HandleCallback", zap.String("data", string(tc.UpdateBotCallbackQuery.Data)))

	ac.State.LockState(logger.Named(logging.LoggerNameLockState))
	defer ac.State.UnlockState(logger.Named(logging.LoggerNameLockState))

	if ac.State.callbackHandler != nil {
		result := ac.State.callbackHandler(string(tc.UpdateBotCallbackQuery.Data))

		if result == nil {
			logger.Warn("CallbackHandler returned nil")
			return
		}

		internalActionHandle(ac, &tc.TelegramUpdateContext, result.Action)

		if !result.NoAnswer {
			tc.AnswerCallbackQuery()
		}

	} else {
		logger.Warn("Missing CallbackHandler")
	}

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}

}

func DefaultHandleMessage[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramContextTextMessage) {
	logger := ac.Loggers.Handler().With(zap.Int64("UpdateID", tc.UpdateID))

	logger.Info("HandleMessage", zap.Any("text", tc.Text))

	ac.State.LockState(logger.Named(logging.LoggerNameLockState))
	defer ac.State.UnlockState(logger.Named(logging.LoggerNameLockState))

	if ac.State.inputHandler != nil {
		ac.State.renderedElements = append(
			ac.State.renderedElements,
			render.NewRenderedUserMessage(tc.Message.ID),
		)

		action := ac.State.inputHandler(tc.Message.Message)

		internalActionHandle(ac, &tc.TelegramUpdateContext, action)

	} else {
		logger.Warn("Missing InputHandler")
	}

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}
}

// Handle
func DefaultHandleActionExternal[S any, C any](ac *ApplicationChat[S, C], tc *telegram.TelegramUpdateContext, action any) {

	actionName := reflection.ReflectStructName(action)
	logger := ac.Loggers.Action().With(zap.String("action", actionName))

	logger.Info("HandleActionExternal")

	ac.State.LockState(logger.Named(logging.LoggerNameLockState))
	defer ac.State.UnlockState(logger.Named(logging.LoggerNameLockState))

	internalActionHandle(ac, tc, action)

	err := ac.App.RenderFunc(tc.Ctx, ac)

	if err != nil {
		logger.Error("Error rendering state", zap.Error(err))
	}

}
