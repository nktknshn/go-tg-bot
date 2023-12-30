package gogotd

import (
	"github.com/go-faster/errors"
	"github.com/gotd/td/tg"
)

func UnpackEditMessage(u tg.UpdatesClass, err error) (*tg.Message, error) {
	if err != nil {
		return nil, err
	}

	var updates []tg.UpdateClass
	switch v := u.(type) {
	case *tg.Updates:
		updates = v.GetUpdates()
	default:
		return nil, errors.Errorf("unexpected type %T", u)
	}

	for _, update := range updates {
		if msgUpdate, ok := update.(*tg.UpdateEditMessage); ok {
			if msg, ok := msgUpdate.Message.(*tg.Message); ok {
				return msg, nil
			}
		}
	}

	return nil, errors.Errorf("no edit message found in %T", u)
}
