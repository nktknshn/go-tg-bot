package dispatcher

import (
	"github.com/nktknshn/go-tg-bot/tgbot/telegram"
)

// handle a certain chat update
type ChatHandler interface {
	HandleUpdate(*telegram.TelegramUpdateContext)
}

// type ApplicationChatHandlerStats struct {
// 	UpdatesAccepted     int
// 	UpdatesSkipped      []string
// 	UpdatesSkippedCount int
// }

// type ApplicationChatHandler[S any, C any] struct {
// 	app     *application.Application[S, C]
// 	appChat *application.ApplicationChat[S, C]
// 	logger  *zap.Logger

// 	// stats
// 	updatesAccepted     int
// 	updatesSkipped      []string
// 	updatesSkippedCount int
// }

// // Creates a new chat handler from an update
// func NewApplicationChatHandler[S any, C any](app *application.Application[S, C], tc *telegram.TelegramUpdateContext) *ApplicationChatHandler[S, C] {

// 	// init new chat
// 	appChat := application.NewApplicationChat[S, C](app, tc)
// 	logger := app.Loggers.ChatHandler(app.Loggers.Base, tc)

// 	return &ApplicationChatHandler[S, C]{
// 		app:     app,
// 		appChat: appChat,
// 		logger:  logger,
// 	}
// }

// // forward update to the application
// func (h *ApplicationChatHandler[S, C]) HandleUpdate(tc *telegram.TelegramUpdateContext) {
// 	h.logger.Debug("HandleUpdate",
// 		zap.Any("UpdateType", tc.Update.UpdateClass.TypeName()),
// 		zap.Any("UpdateID", tc.UpdateID),
// 	)

// 	if tcm, ok := tc.AsTextMessage(); ok && tcm.Message.Message != "" {
// 		h.addAcceptedUpdate(tc.Update)
// 		h.app.HandleMessage(h.appChat, tcm)
// 		return
// 	}

// 	if tccb, ok := tc.AsCallback(); ok {
// 		h.addAcceptedUpdate(tc.Update)
// 		h.app.HandleCallback(h.appChat, tccb)
// 		return
// 	}

// 	// lock happens in the dispatcher's goroutine so we don't need to lock here
// 	h.addSkippedUpdate(tc.Update)

// 	h.logger.Debug("Unkown update (neither message nor callback)")
// }

// func (h *ApplicationChatHandler[S, C]) addAcceptedUpdate(update telegram.BotUpdate) {
// 	h.updatesAccepted++
// }

// func (h *ApplicationChatHandler[S, C]) addSkippedUpdate(update telegram.BotUpdate) {
// 	h.updatesSkipped = append(h.updatesSkipped, update.UpdateClass.TypeName())
// 	h.updatesSkippedCount++

// 	// limit to 30 skipped updates
// 	if len(h.updatesSkipped) > 30 {
// 		h.updatesSkipped = h.updatesSkipped[1:]
// 	}

// 	h.updatesSkippedCount++
// }

// func (h *ApplicationChatHandler[S, C]) Stats() ApplicationChatHandlerStats {
// 	return ApplicationChatHandlerStats{
// 		UpdatesAccepted:     h.updatesAccepted,
// 		UpdatesSkipped:      h.updatesSkipped,
// 		UpdatesSkippedCount: h.updatesSkippedCount,
// 	}
// }
