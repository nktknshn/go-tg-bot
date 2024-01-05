package dispatcher

type ChatsDispatcherStats struct {
	ChatsCount int `json:"chats_count"`
}

func (cd *ChatsDispatcher) Stats() *ChatsDispatcherStats {
	cd.stateLock.Lock()
	defer cd.stateLock.Unlock()
	return &ChatsDispatcherStats{
		ChatsCount: len(cd.chatHandlers),
	}
}
