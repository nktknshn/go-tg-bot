package gogotd

import (
	"fmt"

	"github.com/gotd/td/tg"
)

type UpdatesExtract struct {
	Entities tg.Entities
	Updates  []tg.UpdateClass
}

func ExtractUpdates(updates tg.UpdatesClass) (UpdatesExtract, error) {
	var (
		e    tg.Entities
		upds []tg.UpdateClass
	)

	switch u := updates.(type) {
	case *tg.Updates:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
	case *tg.UpdatesCombined:
		upds = u.Updates
		e.Users = u.MapUsers().NotEmptyToMap()
		chats := u.MapChats()
		e.Chats = chats.ChatToMap()
		e.Channels = chats.ChannelToMap()
	case *tg.UpdateShort:
		upds = []tg.UpdateClass{u.Update}
	default:
		// *UpdateShortMessage
		// *UpdateShortChatMessage
		// *UpdateShortSentMessage
		// *UpdatesTooLong
		return UpdatesExtract{}, fmt.Errorf("unsupported update type %T", updates)
	}

	return UpdatesExtract{
		Entities: e,
		Updates:  upds,
	}, nil
}

func GetChat(u tg.UpdateClass) (*tg.PeerChat, bool) {
	switch u := u.(type) {
	case *tg.UpdateNewMessage:
		println("UpdateNewMessage")
		if m, ok := u.Message.(*tg.Message); ok {
			chat, ok := m.PeerID.(*tg.PeerChat)
			return chat, ok
		}
		return nil, false
	case *tg.UpdateBotCallbackQuery:
		chat, ok := u.Peer.(*tg.PeerChat)
		if !ok {
			return nil, false
		}
		return chat, true
	}
	return nil, false
}

func GetPeerUser(u tg.UpdateClass) (*tg.PeerUser, bool) {
	switch u := u.(type) {
	case *tg.UpdateNewMessage:
		if m, ok := u.Message.(*tg.Message); ok {
			chat, ok := m.PeerID.(*tg.PeerUser)
			return chat, ok
		}
		return nil, false
	case *tg.UpdateBotCallbackQuery:
		chat, ok := u.Peer.(*tg.PeerUser)
		if !ok {
			return nil, false
		}
		return chat, true
	}
	return nil, false
}

func GetUser(entities tg.Entities, u tg.UpdateClass) (*tg.User, bool) {
	peer, ok := GetPeerUser(u)

	if !ok {
		return nil, false
	}

	user, ok := entities.Users[peer.UserID]

	if !ok {
		return nil, false
	}

	return user, true
}
